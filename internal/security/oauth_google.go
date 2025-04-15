package security

import (
	"Forum/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config

func init() {
	// Charger le fichier .env une seule fois
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Erreur de chargement du fichier .env :", err)
	}

	// Lire les valeurs du .env
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	// Initialisation de la configuration OAuth après chargement des variables
	googleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "https://localhost:8443/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

// Redirection vers Google pour authentification
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("randomstate", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "select_account"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Callback après connexion Google
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code de connexion manquant", http.StatusBadRequest)
		return
	}

	// Échanger le code contre un token d'accès
	token, err := googleOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Erreur lors de l'échange du token", http.StatusInternalServerError)
		return
	}

	// Récupérer les infos utilisateur depuis l'API Google
	client := googleOAuthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des infos utilisateur", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Lire la réponse JSON
	var googleUser struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err = json.NewDecoder(resp.Body).Decode(&googleUser)
	if err != nil {
		http.Error(w, "Erreur de lecture des infos utilisateur", http.StatusInternalServerError)
		return
	}

	fmt.Println("Utilisateur Google récupéré :", googleUser.Name, googleUser.Email, googleUser.ID)

	// Vérifier si l'utilisateur existe déjà en base
	user, err := models.GetUserByEmail(googleUser.Email)
	if err != nil {
		fmt.Println("ERREUR SQL lors de la vérification de l'utilisateur :", err)
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	// Si l'utilisateur n'existe pas, le créer
	if user == nil {
		fmt.Println("Utilisateur non trouvé, tentative de création :", googleUser.Name, googleUser.Email, googleUser.ID)
		err = models.CreateGoogleUser(googleUser.Name, googleUser.Email, googleUser.ID)
		if err != nil {
			http.Error(w, "Erreur lors de la création de l'utilisateur", http.StatusInternalServerError)
			return
		}

		// Récupérer l'utilisateur après l'insertion
		user, err = models.GetUserByEmail(googleUser.Email)
		if err != nil || user == nil {
			http.Error(w, "ERREUR: Impossible de récupérer l'utilisateur après insertion", http.StatusInternalServerError)
			return
		}
	}

	// Créer le cookie pour l'utilisateur
	err = CreateCookie(w, r, user.ID, user.Role)
	if err != nil {
		http.Error(w, "Erreur lors de la création du cookie", http.StatusInternalServerError)
		return
	}

	// Redirection via un script JS qui ferme la pop-up et recharge la page principale
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
    <script>
        if (window.opener) {
            window.opener.location.href = "/user/profile"; // Rediriger vers /profile
            window.close(); // Fermer la pop-up
        } else {
            window.location.href = "/user/profile"; // Si pas d'opener, rediriger normalement
        }
    </script>
`)
}
