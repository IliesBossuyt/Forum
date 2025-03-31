package router

import (
	"Forum/internal/database"
	"Forum/internal/handlers"
	"Forum/internal/security"
	"fmt"
	"net/http"
)

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func Router() {
	database.InitDatabase()

	mainRouter := http.NewServeMux()

	// === Middleware de rôle ===
	requireRole := security.RequireRole

	// === Fichiers statiques ===
	mainRouter.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public/static"))))

	// === Guest routes ===
	guestRouter := http.NewServeMux()
	guestRouter.HandleFunc("/", handlers.Accueil)
	guestRouter.HandleFunc("/home", handlers.Home)
	mainRouter.Handle("/", http.StripPrefix("/", guestRouter))

	// === Auth routes ===
	authRouter := http.NewServeMux()
	authRouter.HandleFunc("/register", handlers.Register)
	authRouter.HandleFunc("/login", handlers.Login)
	authRouter.HandleFunc("/logout", handlers.Logout)
	authRouter.HandleFunc("/auth/google/login", security.GoogleLogin)
	authRouter.HandleFunc("/auth/google/callback", security.GoogleCallback)
	authRouter.HandleFunc("/auth/github/login", security.GitHubLogin)
	authRouter.HandleFunc("/auth/github/callback", security.GitHubCallback)
	mainRouter.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	// === User routes ===
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("/profile", handlers.Profile)
	userRouter.HandleFunc("/create-post", handlers.CreatePost)
	userRouter.HandleFunc("/like", handlers.LikePost)
	userRouter.HandleFunc("/edit-post", handlers.EditPost)
	userRouter.HandleFunc("/delete-post", handlers.DeletePost)
	userRouter.HandleFunc("/image/", handlers.GetImage)
	mainRouter.Handle("/user/", requireRole("user")(http.StripPrefix("/user", userRouter)))

	// === Admin routes (Role: Admin) ===
	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc("/dashboard", handlers.DashboardHandler)
	adminRouter.HandleFunc("/change-role", handlers.ChangeUserRole)
	adminRouter.HandleFunc("/toggle-ban", security.ToggleBanUser)
	mainRouter.Handle("/admin/", requireRole("admin")(http.StripPrefix("/admin", adminRouter)))

	go func() {
		// Redirection HTTP vers HTTPS
		err := http.ListenAndServe(":8080", http.HandlerFunc(redirectToHTTPS))
		if err != nil {
			fmt.Println("Erreur redirection HTTP → HTTPS :", err)
		}
	}()

	fmt.Println("Serveur HTTPS lancé sur https://localhost:8443")
	err := http.ListenAndServeTLS(
		":8443",
		"certs/localhost.crt",
		"certs/localhost.key",
		mainRouter,
	)
	
	if err != nil {
		fmt.Println("Erreur serveur HTTPS :", err)
	}

}

