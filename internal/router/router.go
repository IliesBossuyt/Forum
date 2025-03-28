package router

import (
	"fmt"
	"log"
	"net/http"
	"Forum/internal/database"
	"Forum/internal/handlers"
)

func Router() {
	database.InitDatabase()

	// Définir les différentes routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Accueil)
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/profile", handlers.Profile)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.HandleFunc("/home", handlers.Home)
	mux.HandleFunc("/create-post", handlers.CreatePost)
	mux.HandleFunc("/like", handlers.LikePost)

	// Servir les fichiers statiques
	fs := http.FileServer(http.Dir("public/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Redirection HTTP vers HTTPS
	go func() {
		log.Println("Serveur HTTP lancé sur http://localhost:8080 (redirection vers HTTPS)")
		err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://localhost:8443"+r.RequestURI, http.StatusMovedPermanently)
		}))
		if err != nil {
			log.Fatalf("Erreur serveur HTTP : %v", err)
		}
	}()

	// Configuration du serveur HTTPS avec certificat auto-signé
	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
	}

	fmt.Println("Serveur HTTPS lancé sur https://localhost:8443")
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
