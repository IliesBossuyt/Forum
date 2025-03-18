package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Clé secrète pour sécuriser les tokens
var secretKey = []byte("super-secret-key")

// Générer un token sécurisé (UUID + Signature + Rôle)
func GenerateSecureToken(userID, userAgent, role string) (string, error) {
	sessionUUID := uuid.New().String()

	h := hmac.New(sha256.New, secretKey)
	data := sessionUUID + ":" + userAgent + ":" + role
	h.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	token := sessionUUID + ":" + signature + ":" + role

	return token, nil
}

// Vérifier si un token est valide (User-Agent)
func ValidateSecureToken(token, currentUserAgent string) (string, string, bool) {
	parts := splitToken(token, 3)
	if len(parts) != 3 {
		return "", "", false
	}

	sessionUUID, receivedSignature, role := parts[0], parts[1], parts[2]

	h := hmac.New(sha256.New, secretKey)
	data := sessionUUID + ":" + currentUserAgent + ":" + role
	h.Write([]byte(data))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	if hmac.Equal([]byte(receivedSignature), []byte(expectedSignature)) {
		// 🔹 Récupérer `userID` et `role` depuis la session en base
		userID, storedRole, err := GetUserIDFromSession(sessionUUID)
		if err != nil {
			return "", "", false
		}

		// 🔹 Vérifier que le rôle stocké correspond au rôle du token
		if storedRole != role {
			return "", "", false
		}

		return userID, role, true
	}

	return "", "", false
}

// Créer un cookie sécurisé
func CreateCookie(w http.ResponseWriter, r *http.Request, userID, role string) error {
	userAgent := r.UserAgent()

	// Générer le token sécurisé
	token, err := GenerateSecureToken(userID, userAgent, role)
	if err != nil {
		return err
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	// Insérer en base
	sessionUUID := ExtractUUID(token)
	err = CreateSession(sessionUUID, userID, userAgent, role, expirationTime)
	if err != nil {
		return err
	}

	// Stocker le cookie
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

// Supprimer un cookie et la session en base
func DeleteCookie(w http.ResponseWriter, token string) {
	// Supprimer en base
	sessionUUID := ExtractUUID(token)
	DeleteSession(sessionUUID)

	// Supprimer le cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

// Extraire l'UUID depuis un token sécurisé
func ExtractUUID(token string) string {
	parts := splitToken(token, 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// Séparer un token en parties
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
