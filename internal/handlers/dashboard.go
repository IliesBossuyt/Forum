package handlers

import (
	"Forum/internal/database"
	"Forum/internal/models"
	"Forum/internal/security"
	"html/template"
	"net/http"
)

// Données nécessaires pour le dashboard
type DashboardData struct {
	Users           []models.User
	Reports         []models.Report
	CommentReports  []models.CommentReport
	WarnCounts      map[string]int
	CurrentUserRole string
}

// Affiche le dashboard
func Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Récupère tous les utilisateurs
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

		// Récupère les signalements de commentaires
		commentReports, err := models.GetAllCommentReports()
		if err != nil {
			http.Error(w, "Erreur serveur (signalements commentaires)", http.StatusInternalServerError)
			return
		}

		// Récupère les signalements de posts
		reports, err := models.GetAllReports()
		if err != nil {
			http.Error(w, "Erreur serveur (signalements)", http.StatusInternalServerError)
			return
		}

		// Récupère les posts pour lier les auteurs
		posts, err := models.GetAllPosts()
		if err != nil {
			http.Error(w, "Erreur serveur (posts)", http.StatusInternalServerError)
			return
		}

		// Crée une map pour lier les posts à leurs auteurs
		postAuthorMap := make(map[int]string)
		for _, post := range posts {
			postAuthorMap[post.ID] = post.UserID
		}

		// Ajoute l'ID de l'auteur à chaque signalement
		for i := range reports {
			if authorID, ok := postAuthorMap[reports[i].PostID]; ok {
				reports[i].PostAuthorID = authorID
			}
		}

		// Récupère et compte les avertissements
		warns, err := models.GetAllWarns()
		if err != nil {
			http.Error(w, "Erreur récupération des warns", http.StatusInternalServerError)
			return
		}

		warnCounts := make(map[string]int)
		for _, warn := range warns {
			warnCounts[warn.UserID]++
		}

		// Charge et exécute le template
		tmpl, err := template.ParseFiles("../public/template/dashboard.html")
		if err != nil {
			http.Error(w, "Erreur de chargement du template", http.StatusInternalServerError)
			return
		}

		role := r.Context().Value(security.ContextRoleKey).(string)
		data := DashboardData{
			Users:           users,
			Reports:         reports,
			CommentReports:  commentReports,
			WarnCounts:      warnCounts,
			CurrentUserRole: role,
		}

		tmpl.Execute(w, data)
	}
}
