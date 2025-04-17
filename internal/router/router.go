package router

import (
	"Forum/internal/database"
	"Forum/internal/handlers"
	"Forum/internal/models"
	"Forum/internal/security"
	"fmt"
	"net/http"
	"time"
)

// Redirige les requêtes HTTP vers HTTPS
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

// Configure et lance le serveur web
func Router() {
	// Initialise la base de données
	database.InitDatabase()

	// Crée le routeur principal
	routeManager := http.NewServeMux()

	// Middleware pour vérifier les rôles
	requireRole := security.RequireRole

	// Routes pour les fichiers statiques
	routeManager.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../public/static"))))

	// Routes accessibles aux invités
	guestRouter := http.NewServeMux()
	guestRouter.HandleFunc("/home", handlers.Home)
	guestRouter.HandleFunc("/image/", handlers.GetImage)
	routeManager.Handle("/entry/", requireRole("guest", "user", "admin", "moderator")(
		http.StripPrefix("/entry", handlers.WithNotFoundFallback(guestRouter)),
	))

	// Routes d'authentification
	authRouter := http.NewServeMux()
	authRouter.Handle("/register", security.RateLimitRegisterByIP(http.HandlerFunc(handlers.Register)))
	authRouter.Handle("/login", security.RateLimitLoginByIP(security.RateLimitLoginByIdentifier(http.HandlerFunc(handlers.Login))))
	authRouter.HandleFunc("/logout", handlers.Logout)
	authRouter.HandleFunc("/unauthorized", handlers.UnauthorizedHandler)
	authRouter.HandleFunc("/google/login", security.GoogleLogin)
	authRouter.HandleFunc("/google/callback", security.GoogleCallback)
	authRouter.HandleFunc("/github/login", security.GitHubLogin)
	authRouter.HandleFunc("/github/callback", security.GitHubCallback)
	routeManager.Handle("/auth/", http.StripPrefix("/auth", handlers.WithNotFoundFallback(authRouter)))

	// Routes utilisateur
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("/profile/", handlers.Profile)
	userRouter.Handle("/create-post", security.RateLimitCreatePost(http.HandlerFunc(handlers.CreatePost)))
	userRouter.HandleFunc("/like", handlers.LikePost)
	userRouter.HandleFunc("/edit-post", handlers.EditPost)
	userRouter.HandleFunc("/delete-post", handlers.DeletePost)
	userRouter.HandleFunc("/report", handlers.ReportPost)
	userRouter.HandleFunc("/comment", handlers.PostComment)
	userRouter.HandleFunc("/like-comment", handlers.LikeComment)
	userRouter.HandleFunc("/delete-comment", handlers.DeleteComment)
	userRouter.HandleFunc("/edit-comment", handlers.EditComment)
	userRouter.HandleFunc("/report-comment", handlers.ReportComment)
	userRouter.HandleFunc("/notifications", handlers.GetNotifications)
	userRouter.HandleFunc("/notifications/mark-read", handlers.MarkNotificationsRead)
	userRouter.HandleFunc("/notifications/delete-all", handlers.DeleteAllNotifications)
	routeManager.Handle("/user/", requireRole("user", "admin", "moderator")(http.StripPrefix("/user", userRouter)))

	// Routes administrateur
	adminRouter := http.NewServeMux()
	// Routes accessibles aux administrateurs et modérateurs
	adminRouter.HandleFunc("/dashboard", handlers.Dashboard)
	adminRouter.HandleFunc("/delete-report-post", handlers.DeleteReportPost)
	adminRouter.HandleFunc("/delete-report-comment", handlers.DeleteReportComment)
	adminRouter.HandleFunc("/add-warn", handlers.AddWarn)
	adminRouter.HandleFunc("/delete-comment", handlers.DeleteComment)
	adminRouter.HandleFunc("/delete-post", handlers.DeletePost)
	adminRouter.HandleFunc("/warns", handlers.GetUserWarns)

	// Routes sensibles (administrateurs uniquement)
	adminSecure := http.NewServeMux()
	adminSecure.HandleFunc("/change-role", handlers.ChangeUserRole)
	adminSecure.HandleFunc("/toggle-ban", security.ToggleBanUser)
	adminSecure.HandleFunc("/delete-warn", handlers.DeleteWarn)

	// Attache les routes avec les bons niveaux d'accès
	routeManager.Handle("/admin/", requireRole("admin", "moderator")(http.StripPrefix("/admin", handlers.WithNotFoundFallback(adminRouter))))
	routeManager.Handle("/admin/secure/", requireRole("admin")(http.StripPrefix("/admin/secure", handlers.WithNotFoundFallback(adminSecure))))

	// Handler final avec gestion des erreurs 404
	var secureHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler, pattern := routeManager.Handler(r)
		if pattern == "" {
			handlers.NotFoundHandler(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})

	// Applique le rate limiting global
	secureHandler = security.RateLimitGlobal(secureHandler)

	// Nettoie périodiquement les sessions expirées
	go func() {
		for {
			models.CleanExpiredSessions()
			time.Sleep(24 * time.Hour)
		}
	}()

	// Redirige le trafic HTTP vers HTTPS
	go func() {
		err := http.ListenAndServe(":8080", http.HandlerFunc(redirectToHTTPS))
		if err != nil {
			fmt.Println("Erreur redirection HTTP → HTTPS :", err)
		}
	}()

	// Lance le serveur HTTPS
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
