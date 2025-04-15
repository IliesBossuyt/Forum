package security

import (
	"Forum/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type contextKey string

const (
	ContextUserIDKey contextKey = "userID"
	ContextRoleKey   contextKey = "role"
)

// Middleware pour la vérification des roles, des sessions et du ban
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID, role string
			role = "guest" // Valeur par défaut pour les visiteurs

			cookie, err := r.Cookie("session")
			if err == nil {
				// Vérifier le token (signature et User-Agent)
				uid, userRole, valid := ValidateSecureToken(cookie.Value, r.UserAgent())
				if valid {
					userID = uid
					role = userRole

					// Vérifier la session en base de données
					sessionUUID := ExtractUUID(cookie.Value)
					storedUserID, storedRole, expiresAt, err := models.GetUserIDFromSession(sessionUUID)
					if err != nil || storedUserID != userID || storedRole != role || expiresAt.Before(time.Now()) {
						http.Error(w, "Session invalide ou expirée", http.StatusUnauthorized)
						return
					}

					// Vérifier si l'utilisateur est banni
					user, err := models.GetUserByID(userID)
					if err != nil || user == nil || user.Banned {
						http.Error(w, "Utilisateur invalide ou banni", http.StatusUnauthorized)
						return
					}
				}
			}

			// Si rôle non autorisé → redirection
			if !contains(allowedRoles, role) {
                http.Redirect(w, r, "/auth/unauthorized", http.StatusSeeOther)
				return
			}

			// Ajout des infos au contexte
			// Injecter l'utilisateur dans le contexte
			ctx := context.WithValue(r.Context(), ContextUserIDKey, userID)
			ctx = context.WithValue(ctx, ContextRoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func contains(roles []string, target string) bool {
	for _, role := range roles {
		if role == target {
			return true
		}
	}
	return false
}

type Bucket struct {
	tokens      int
	lastUpdated time.Time
	mu          sync.Mutex
}

var buckets = make(map[string]*Bucket)
var bucketsMu sync.Mutex

// NewRateLimitMiddleware crée un middleware personnalisé avec un générateur de clé
func NewRateLimitMiddleware(rate int, per time.Duration, keyFunc func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)

			if key == "" {
				// Pas de clé → pas de rate limit
				next.ServeHTTP(w, r)
				return
			}

			if !allowRequest(key, rate, per) {
				http.Error(w, "Trop de requêtes, réessayez plus tard", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func allowRequest(key string, rate int, per time.Duration) bool {
	bucketsMu.Lock()
	bucket, exists := buckets[key]
	if !exists {
		bucket = &Bucket{tokens: rate, lastUpdated: time.Now()}
		buckets[key] = bucket
	}
	bucketsMu.Unlock()

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastUpdated)
	newTokens := int(elapsed / per)

	if newTokens > 0 {
		bucket.tokens = min(bucket.tokens+newTokens, rate)
		bucket.lastUpdated = now
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// === Fonctions utilitaires pour générer des clés ===

func GetIP(r *http.Request) string {
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip, _, _ = net.SplitHostPort(ip)
	}
	return ip
}

func EmailFromRequest(r *http.Request) string {
	r.ParseForm()
	email := r.FormValue("identifier")
	if email == "" {
		email = r.FormValue("email")
	}
	return strings.ToLower(email)
}

func IdentifierKey(r *http.Request) string {
	if r.Method != http.MethodPost {
		return ""
	}

	var bodyCopy bytes.Buffer
	tee := io.TeeReader(r.Body, &bodyCopy)
	r.Body = io.NopCloser(&bodyCopy)

	var data struct {
		Identifier string `json:"identifier"`
	}

	err := json.NewDecoder(tee).Decode(&data)
	if err != nil {
		return "unknown"
	}

	return strings.ToLower(data.Identifier)
}

func UserIDFromContext(r *http.Request) string {
	id, _ := r.Context().Value(ContextUserIDKey).(string)
	return id
}

// Middleware pour le login
var RateLimitLoginByIP = NewRateLimitMiddleware(10, time.Minute, func(r *http.Request) string {
	if r.Method != http.MethodPost {
		return "" // Ne limite pas les requêtes non POST
	}
	return "login-ip:" + GetIP(r)
})

var RateLimitLoginByIdentifier = NewRateLimitMiddleware(5, time.Minute, func(r *http.Request) string {
	if r.Method != http.MethodPost {
		return ""
	}
	return "login-id:" + IdentifierKey(r)
})

// Middleware pour le register
var RateLimitRegisterByIP = NewRateLimitMiddleware(3, time.Minute, func(r *http.Request) string {
	if r.Method != http.MethodPost {
		return ""
	}
	return "register-ip:" + GetIP(r)
})

// Middleware pour la création de post
var RateLimitCreatePost = NewRateLimitMiddleware(5, time.Minute, func(r *http.Request) string {
	return "createpost-userid:" + UserIDFromContext(r)
})

// Middleware global
var RateLimitGlobal = NewRateLimitMiddleware(300, time.Second, func(r *http.Request) string {
	return "global:" + GetIP(r)
})
