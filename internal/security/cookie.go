package security

import (
	"Forum/internal/models"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Clé secrète pour la signature des tokens
var secretKey = []byte("super-secret-key")

// Génère un token sécurisé avec UUID, signature et rôle
func GenerateSecureToken(userID, userAgent, role string) (string, error) {
	// Génère un UUID unique pour la session
	sessionUUID := uuid.New().String()

	// Crée une signature HMAC
	h := hmac.New(sha256.New, secretKey)
	data := sessionUUID + ":" + userAgent + ":" + role
	h.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Combine les éléments pour former le token
	token := sessionUUID + ":" + signature + ":" + role

	return token, nil
}

// Vérifie la validité d'un token et retourne l'ID utilisateur et le rôle
func ValidateSecureToken(token, currentUserAgent string) (string, string, bool) {
	// Découpe le token en ses composants
	parts := splitToken(token, 3)
	if len(parts) != 3 {
		return "", "", false
	}

	sessionUUID, receivedSignature, role := parts[0], parts[1], parts[2]

	// Vérifie la signature
	h := hmac.New(sha256.New, secretKey)
	data := sessionUUID + ":" + currentUserAgent + ":" + role
	h.Write([]byte(data))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	if hmac.Equal([]byte(receivedSignature), []byte(expectedSignature)) {
		// Récupère les informations de session
		userID, storedRole, expiresAt, err := models.GetUserIDFromSession(sessionUUID)
		if err != nil {
			return "", "", false
		}

		// Vérifie le rôle et la date d'expiration
		if storedRole != role || expiresAt.Before(time.Now()) {
			return "", "", false
		}

		return userID, role, true
	}

	return "", "", false
}

// Crée un cookie sécurisé pour une session
func CreateCookie(w http.ResponseWriter, r *http.Request, userID, role string) error {
	userAgent := r.UserAgent()

	// Vérifie l'existence de l'utilisateur
	user, err := models.GetUserByID(userID)
	if err != nil || user == nil {
		// Gère l'erreur
	}

	// Génère le token sécurisé
	token, err := GenerateSecureToken(userID, userAgent, role)
	if err != nil {
		return err
	}

	// Définit la date d'expiration
	expirationTime := time.Now().Add(24 * time.Hour)

	// Crée la session en base de données
	sessionUUID := ExtractUUID(token)
	err = models.CreateSession(sessionUUID, userID, userAgent, role, expirationTime)
	if err != nil {
		return err
	}

	// Définit le cookie HTTP
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	return nil
}

// Supprime un cookie et sa session associée
func DeleteCookie(w http.ResponseWriter, token string) {
	// Supprime la session de la base de données
	sessionUUID := ExtractUUID(token)
	models.DeleteSession(sessionUUID)

	// Supprime le cookie du navigateur
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

// Extrait l'UUID d'un token sécurisé
func ExtractUUID(token string) string {
	parts := splitToken(token, 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// Découpe une chaîne en parties selon un séparateur
func splitToken(s string, n int) []string {
	parts := make([]string, 0, n)
	start := 0
	for i := 0; i < len(s) && len(parts) < n-1; i++ {
		if s[i] == ':' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}
