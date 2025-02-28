package server

import (
	"fmt"
	"net/http"
)

func Router(forum *Forum) {
	// Définir les différentes routes
	http.HandleFunc("/", forum.Accueil)

	fs := http.FileServer(http.Dir("front/./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/forum", forum.Accueil)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Serveur lancé sur http://localhost:8080/")
	// On lance le serveur local sur le port 8080
}
