package server

import (
	"fmt"
	"net/http"
)

func Router(forum *Forum) {
	forum.Init()
	// Définir les différentes routes
	http.HandleFunc("/", Accueil)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/profile", Profile)
	http.HandleFunc("/logout", Logout)


	fs := http.FileServer(http.Dir("front/./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/forum", Accueil)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Serveur lancé sur http://localhost:8080/")
	// On lance le serveur local sur le port 8080
}
