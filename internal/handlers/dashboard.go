package handlers

import (
	"Forum/internal/database"
	"Forum/internal/models"
	"html/template"
	"net/http"
)

// Structure de données pour affichage
type DashboardData struct {
	Users   []models.User
	Reports []models.Report
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Récupération des utilisateurs
		rows, err := database.DB.Query("SELECT id, username, role, banned FROM users")
		if err != nil {
			http.Error(w, "Erreur serveur (utilisateurs)", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var u models.User
			if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.Banned); err != nil {
				http.Error(w, "Erreur lecture utilisateurs", http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}

		// Récupération des signalements
		reports, err := models.GetAllReports()
		if err != nil {
			http.Error(w, "Erreur serveur (signalements)", http.StatusInternalServerError)
			return
		}

		// Récupération des posts (pour obtenir l'auteur via PostID)
		posts, err := models.GetAllPosts()
		if err != nil {
			http.Error(w, "Erreur serveur (posts)", http.StatusInternalServerError)
			return
		}

		// Création d'une map PostID -> UserID
		postAuthorMap := make(map[int]string)
		for _, post := range posts {
			postAuthorMap[post.ID] = post.UserID
		}

		// Injection dynamique du PostAuthorID dans chaque report
		for i := range reports {
			if authorID, ok := postAuthorMap[reports[i].PostID]; ok {
				reports[i].PostAuthorID = authorID
			}
		}

		// Template
		tmpl, err := template.ParseFiles("../public/template/dashboard.html")
		if err != nil {
			http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
			return
		}

		data := DashboardData{
			Users:   users,
			Reports: reports,
		}

		tmpl.Execute(w, data)
	}
}
