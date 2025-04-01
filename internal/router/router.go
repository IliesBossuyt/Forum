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
	// Initialisation de la base de données
	database.InitDatabase()

	// Routeur principal de configuration des routes
	routeManager := http.NewServeMux()

	// Middleware
	requireRole := security.RequireRole

	// Fichiers statiques
	routeManager.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../public/static"))))

	// Guest routes
	guestRouter := http.NewServeMux()
	guestRouter.HandleFunc("/", handlers.Accueil)
	guestRouter.HandleFunc("/home", handlers.Home)
	guestRouter.HandleFunc("/image/", handlers.GetImage)
	routeManager.Handle("/entry/", requireRole("guest", "user", "admin", "moderator")(http.StripPrefix("/entry", guestRouter)))

	// Auth routes
	authRouter := http.NewServeMux()
	authRouter.Handle("/register", security.RateLimitRegisterByIP(http.HandlerFunc(handlers.Register)))
	authRouter.Handle("/login", security.RateLimitLoginByIP(security.RateLimitLoginByIdentifier(http.HandlerFunc(handlers.Login))))
	authRouter.HandleFunc("/logout", handlers.Logout)
	authRouter.HandleFunc("/unauthorized", handlers.UnauthorizedHandler)
	authRouter.HandleFunc("/google/login", security.GoogleLogin)
	authRouter.HandleFunc("/google/callback", security.GoogleCallback)
	authRouter.HandleFunc("/github/login", security.GitHubLogin)
	authRouter.HandleFunc("/github/callback", security.GitHubCallback)
	routeManager.Handle("/auth/", http.StripPrefix("/auth", authRouter))

	// User routes
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("/profile", handlers.Profile)
	userRouter.Handle("/create-post", security.RateLimitCreatePost(http.HandlerFunc(handlers.CreatePost)))
	userRouter.HandleFunc("/like", handlers.LikePost)
	userRouter.HandleFunc("/edit-post", handlers.EditPost)
	userRouter.HandleFunc("/delete-post", handlers.DeletePost)
	routeManager.Handle("/user/", requireRole("user", "admin", "moderator")(http.StripPrefix("/user", userRouter)))

	// Admin routes
	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc("/dashboard", handlers.DashboardHandler)
	adminRouter.HandleFunc("/change-role", handlers.ChangeUserRole)
	adminRouter.HandleFunc("/toggle-ban", security.ToggleBanUser)
	routeManager.Handle("/admin/", requireRole("admin")(http.StripPrefix("/admin", adminRouter)))

	// Handler final avec fallback 404
	var secureHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, pattern := routeManager.Handler(r)
		if pattern == "" {
			handlers.NotFoundHandler(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
	
	// Appliquer le rate limit global
	secureHandler = security.RateLimitGlobal(secureHandler)

	// Redirection HTTP → HTTPS
	go func() {
		err := http.ListenAndServe(":8080", http.HandlerFunc(redirectToHTTPS))
		if err != nil {
			fmt.Println("Erreur redirection HTTP → HTTPS :", err)
		}
	}()

	// Lancement du serveur HTTPS
	fmt.Println("Serveur HTTPS lancé sur https://localhost:8443")
	err := http.ListenAndServeTLS(
		":8443",
		"certs/localhost.crt",
		"certs/localhost.key",
		secureHandler,
	)

	if err != nil {
		fmt.Println("Erreur serveur HTTPS :", err)
	}
}
