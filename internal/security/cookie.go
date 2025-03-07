package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net"
	"net/http"
	"time"
)

// Clé secrète pour sécuriser les tokens (À NE JAMAIS PARTAGER)
var secretKey = []byte("super-secret-key")

// Générer un token sécurisé basé sur un UUID + User-Agent + IP
func GenerateSecureToken(userID, userAgent, userIP string) (string, error) {
	h := hmac.New(sha256.New, secretKey)
	data := userID + ":" + userAgent + ":" + userIP
	h.Write([]byte(data))
	signature := h.Sum(nil)

	// Encoder UUID + signature en base64
	token := base64.URLEncoding.EncodeToString([]byte(userID + ":" + base64.URLEncoding.EncodeToString(signature)))
	return token, nil
}

// Vérifier si le token est valide
func ValidateSecureToken(token, currentUserAgent, currentUserIP string) (string, bool) {
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", false
	}

	// Séparer UUID et signature
	parts := string(decoded)
	split := splitToken(parts, 2) // Séparer UUID et signature
	if len(split) != 2 {
		return "", false
	}

	userID, receivedSignature := split[0], split[1]

	// Recalculer la signature attendue
	h := hmac.New(sha256.New, secretKey)
	data := userID + ":" + currentUserAgent + ":" + currentUserIP
	h.Write([]byte(data))
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Vérifier si la signature est correcte
	if hmac.Equal([]byte(receivedSignature), []byte(expectedSignature)) {
		return userID, true
	}
	return "", false
}

// Séparer les parties du token
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

// Fonction pour extraire uniquement l'IP (sans le port)
func ExtractIP(remoteAddr string) string {
	// Si l'adresse contient ":", il y a un port → on le coupe
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr // Retourne l'adresse originale si erreur
	}
	return host
}

func CreateCookie(w http.ResponseWriter, r *http.Request, userID string) {
	// Récupérer l'empreinte du navigateur (User-Agent + IP)
	userAgent := r.UserAgent()
	userIP := ExtractIP(r.RemoteAddr)

	// Générer un token sécurisé lié à cet appareil
	token, err := GenerateSecureToken(userID, userAgent, userIP)
	if err != nil {
		http.Error(w, "Erreur création token sécurisé", http.StatusInternalServerError)
		return
	}

	// 🔹 Créer un cookie sécurisé
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func DeleteCookie(w http.ResponseWriter) {
	// Supprimer le cookie de session
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
