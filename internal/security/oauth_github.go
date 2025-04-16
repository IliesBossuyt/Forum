package security

import (
	"Forum/internal/models"
	"encoding/json"
	"fmt"
	"strings"

	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubOAuthConfig *oauth2.Config

func init() {
	// Charger le fichier .env
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Erreur de chargement du fichier .env :", err)
	}

	// Lire les valeurs du .env
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	// Initialisation de la configuration OAuth
	githubOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "https://localhost:8443/auth/github/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

// Redirection vers GitHub pour authentification
func GitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOAuthConfig.AuthCodeURL("randomstate", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Structure pour stocker la réponse de GitHub
type GitHubUser struct {
	Login string `json:"login"` // Pseudo GitHub
	Name  string `json:"name"`  // Nom complet (peut être vide)
	Email string `json:"email"` // Email
	ID    int    `json:"id"`    // ID unique GitHub
}

// Callback après connexion GitHub
func GitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code de connexion manquant", http.StatusBadRequest)
		return
	}

	// Échanger le code contre un token d'accès
	token, err := githubOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Erreur lors de l'échange du token GitHub", http.StatusInternalServerError)
		return
	}

	// Récupérer les infos utilisateur depuis l'API GitHub
	client := githubOAuthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des infos utilisateur GitHub", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Lire la réponse JSON
	var githubUser GitHubUser
	err = json.NewDecoder(resp.Body).Decode(&githubUser)
	if err != nil {
		http.Error(w, "Erreur de lecture des infos utilisateur GitHub", http.StatusInternalServerError)
		return
	}

	// Si `name` est vide, utiliser `login` comme nom
	username := githubUser.Name
	if username == "" {
		username = githubUser.Login
	}

	// Nettoyer le username (supprimer les espaces)
	username = strings.ReplaceAll(username, " ", "")

	// Vérifier si l'email est vide et récupérer l'email principal si nécessaire
	if githubUser.Email == "" {
		respEmails, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des emails GitHub", http.StatusInternalServerError)
			return
		}
		defer respEmails.Body.Close()

		// Lire les emails disponibles
		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		err = json.NewDecoder(respEmails.Body).Decode(&emails)
		if err != nil {
			http.Error(w, "Erreur de lecture des emails GitHub", http.StatusInternalServerError)
			return
		}

		// Sélectionner l'email principal et vérifié
		for _, e := range emails {
			if e.Primary && e.Verified {
				githubUser.Email = e.Email
				break
			}
		}
	}

	// Vérifier si l'utilisateur existe déjà en base
	user, err := models.GetUserByEmail(githubUser.Email)
	if err != nil {
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	// Si l'utilisateur n'existe pas, le créer
	if user == nil {
		err = models.CreateGitHubUser(username, githubUser.Email, fmt.Sprint(githubUser.ID))
		if err != nil {
			http.Error(w, "Erreur lors de la création de l'utilisateur", http.StatusInternalServerError)
			return
		}

		// Récupérer l'utilisateur après l'insertion
		user, err = models.GetUserByEmail(githubUser.Email)
		if err != nil || user == nil {
			return
		}
	}

	// Créer le cookie et rediriger
	err = CreateCookie(w, r, user.ID, user.Role)
	if err != nil {
		http.Error(w, "Erreur lors de la création du cookie", http.StatusInternalServerError)
		return
	}

	// Redirection vers le profil public du user
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<script>
			if (window.opener) {
				window.opener.location.href = "/entry/home";
				window.close();
			} else {
				window.location.href = "/entry/home";
			}
		</script>
	`,)
}
