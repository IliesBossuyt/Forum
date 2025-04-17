package security

import (
	"Forum/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Configuration OAuth pour Google
var googleOAuthConfig *oauth2.Config

// Initialisation de la configuration OAuth au démarrage
func init() {
	// Charge les variables d'environnement depuis le fichier .env
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Erreur de chargement du fichier .env :", err)
	}

	// Récupère les identifiants OAuth depuis les variables d'environnement
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	// Configure les paramètres OAuth pour Google
	googleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "https://localhost:8443/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

// Redirige l'utilisateur vers la page d'authentification Google
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Configure la redirection avec sélection du compte
	url := googleOAuthConfig.AuthCodeURL("randomstate", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "select_account"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Gère le callback après l'authentification Google
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Récupère le code d'autorisation
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code de connexion manquant", http.StatusBadRequest)
		return
	}

	// Échange le code contre un token d'accès
	token, err := googleOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Erreur lors de l'échange du token", http.StatusInternalServerError)
		return
	}

	// Récupère les informations de l'utilisateur depuis l'API Google
	client := googleOAuthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des infos utilisateur", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Structure pour stocker les informations de l'utilisateur Google
	var googleUser struct {
		ID    string `json:"id"`    // ID unique Google
		Name  string `json:"name"`  // Nom complet
		Email string `json:"email"` // Email principal
	}

	// Décode les informations de l'utilisateur
	err = json.NewDecoder(resp.Body).Decode(&googleUser)
	if err != nil {
		http.Error(w, "Erreur de lecture des infos utilisateur", http.StatusInternalServerError)
		return
	}

	// Vérifie si l'utilisateur existe déjà
	user, err := models.GetUserByEmail(googleUser.Email)
	if err != nil {
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	// Crée un nouvel utilisateur si nécessaire
	if user == nil {
		// Nettoie le nom d'utilisateur (supprime les espaces)
		cleanName := strings.ReplaceAll(googleUser.Name, " ", "")
		err = models.CreateGoogleUser(cleanName, googleUser.Email, googleUser.ID)
		if err != nil {
			http.Error(w, "Erreur lors de la création de l'utilisateur", http.StatusInternalServerError)
			return
		}

		// Récupère l'utilisateur après création
		user, err = models.GetUserByEmail(googleUser.Email)
		if err != nil || user == nil {
			http.Error(w, "ERREUR: Impossible de récupérer l'utilisateur après insertion", http.StatusInternalServerError)
			return
		}
	}

	// Crée la session et redirige
	err = CreateCookie(w, r, user.ID, user.Role)
	if err != nil {
		http.Error(w, "Erreur lors de la création du cookie", http.StatusInternalServerError)
		return
	}

	// Redirige vers la page d'accueil avec gestion de la pop-up
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
    <script>
        if (window.opener) {
            window.opener.location.href = "/entry/home";
            window.close(); // Ferme la pop-up
        } else {
            window.location.href = "/entry/home";
        }
    </script>
    `)
}
