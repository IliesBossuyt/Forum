package router

import (
	"fmt"
	"net/http"

	"Forum/internal/database"
	"Forum/internal/handlers"
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
	

	fs := http.FileServer(http.Dir("public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
	fmt.Println("Serveur lancé sur http://localhost:8080/")
	// On lance le serveur local sur le port 8080
}
