package router

import (
	"fmt"
	"net/http"

	"Forum/internal/database"
	"Forum/internal/handlers"
	"Forum/internal/security"
)

func Router() {
	database.InitDatabase()
	// Définir les différentes routes
	http.HandleFunc("/", handlers.Accueil)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/profile", handlers.Profile)
	http.HandleFunc("/logout", handlers.Logout)
	http.HandleFunc("/home", handlers.Home)
	http.HandleFunc("/create-post", handlers.CreatePost)
	http.HandleFunc("/like", handlers.LikePost)
	http.HandleFunc("/edit-post", handlers.EditPost)
	http.HandleFunc("/image/{id}", handlers.GetImage)
	http.HandleFunc("/dashboard", security.AdminOnly(handlers.DashboardHandler))
	http.HandleFunc("/change-role", security.AdminOnly(handlers.ChangeUserRole))
	http.HandleFunc("/delete-post", handlers.DeletePost)
	http.HandleFunc("/auth/google/login", security.GoogleLogin)
	http.HandleFunc("/auth/google/callback", security.GoogleCallback)
	http.HandleFunc("/auth/github/login", security.GitHubLogin)
	http.HandleFunc("/auth/github/callback", security.GitHubCallback)

	fs := http.FileServer(http.Dir("../public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
	fmt.Println("Serveur lancé sur http://localhost:8080/")
	// On lance le serveur local sur le port 8080
}
