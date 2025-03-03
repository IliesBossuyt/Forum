package router

import (
	"fmt"
	"net/http"

	"Forum/internal/handlers"
	"Forum/internal/database"
)

func Router() {
	database.InitDatabase()
	// Définir les différentes routes
	http.HandleFunc("/", handlers.Accueil)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/profile", handlers.Profile)
	http.HandleFunc("/logout", handlers.Logout)


	fs := http.FileServer(http.Dir("front/./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/forum", handlers.Accueil)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Serveur lancé sur http://localhost:8080/")
	// On lance le serveur local sur le port 8080
}
