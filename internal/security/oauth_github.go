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

// Configuration OAuth pour GitHub
var githubOAuthConfig *oauth2.Config

// Initialisation de la configuration OAuth au démarrage
func init() {
	// Charge les variables d'environnement depuis le fichier .env
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Erreur de chargement du fichier .env :", err)
	}

	// Récupère les identifiants OAuth depuis les variables d'environnement
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	// Configure les paramètres OAuth pour GitHub
	githubOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "https://localhost:8443/auth/github/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

// Redirige l'utilisateur vers la page d'authentification GitHub
func GitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOAuthConfig.AuthCodeURL("randomstate", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Structure pour stocker les informations de l'utilisateur GitHub
type GitHubUser struct {
	Login string `json:"login"` // Nom d'utilisateur GitHub
	Name  string `json:"name"`  // Nom complet
	Email string `json:"email"` // Email principal
	ID    int    `json:"id"`    // ID unique GitHub
}

// Gère le callback après l'authentification GitHub
func GitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Récupère le code d'autorisation
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code de connexion manquant", http.StatusBadRequest)
		return
	}

	// Échange le code contre un token d'accès
	token, err := githubOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Erreur lors de l'échange du token GitHub", http.StatusInternalServerError)
		return
	}

	// Récupère les informations de l'utilisateur depuis l'API GitHub
	client := githubOAuthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des infos utilisateur GitHub", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Décode les informations de l'utilisateur
	var githubUser GitHubUser
	err = json.NewDecoder(resp.Body).Decode(&githubUser)
	if err != nil {
		http.Error(w, "Erreur de lecture des infos utilisateur GitHub", http.StatusInternalServerError)
		return
	}

	// Utilise le login comme nom d'utilisateur si le nom est vide
	username := githubUser.Name
	if username == "" {
		username = githubUser.Login
	}

	// Nettoie le nom d'utilisateur (supprime les espaces)
	username = strings.ReplaceAll(username, " ", "")

	// Récupère l'email principal si nécessaire
	if githubUser.Email == "" {
		respEmails, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des emails GitHub", http.StatusInternalServerError)
			return
		}
		defer respEmails.Body.Close()

		// Structure pour les emails GitHub
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

		// Sélectionne l'email principal vérifié
		for _, e := range emails {
			if e.Primary && e.Verified {
				githubUser.Email = e.Email
				break
			}
		}
	}

	// Vérifie si l'utilisateur existe déjà
	user, err := models.GetUserByEmail(githubUser.Email)
	if err != nil {
		http.Error(w, "Erreur interne", http.StatusInternalServerError)
		return
	}

	// Crée un nouvel utilisateur si nécessaire
	if user == nil {
		err = models.CreateGitHubUser(username, githubUser.Email, fmt.Sprint(githubUser.ID))
		if err != nil {
			http.Error(w, "Erreur lors de la création de l'utilisateur", http.StatusInternalServerError)
			return
		}

		// Récupère l'utilisateur après création
		user, err = models.GetUserByEmail(githubUser.Email)
		if err != nil || user == nil {
			return
		}
	}

	// Crée la session et redirige
	err = CreateCookie(w, r, user.ID, user.Role)
	if err != nil {
		http.Error(w, "Erreur lors de la création du cookie", http.StatusInternalServerError)
		return
	}

	// Redirige vers la page d'accueil
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
	`)
}
