package handlers

import (
	"Forum/internal/database"
	"html/template"
	"net/http"
)

// Structure User pour l'affichage
type User struct {
	ID       string
	Username string
	Role     string
	Banned   bool
}

// Handler pour afficher le tableau de bord admin
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer tous les utilisateurs
	rows, err := database.DB.Query("SELECT id, username, role, banned FROM users")
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.Banned); err != nil {
			http.Error(w, "Erreur de récupération", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	// Charger et exécuter le template HTML
	tmpl, err := template.ParseFiles("../public/template/dashboard.html")
	if err != nil {
		http.Error(w, "Erreur chargement template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, struct{ Users []User }{Users: users})
}
