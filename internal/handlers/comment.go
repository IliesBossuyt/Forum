package handlers

import (
    "Forum/internal/models"
    "Forum/internal/security"
    "log"
    "net/http"
    "strconv"
)

// Handler inutilisé vu que tout passe par la page home.html
func ShowPostWithComments(w http.ResponseWriter, r *http.Request) {
    // Plutôt que de charger une page post.html qui n'existe pas,
    // on redirige vers la home qui affiche déjà tous les posts et leurs commentaires
    http.Redirect(w, r, "/entry/home", http.StatusSeeOther)
}

func AddComment(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Vérifie que l'utilisateur est bien connecté via le cookie
    cookie, err := r.Cookie("session")
    if err != nil {
        http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
        return
    }

    userID, _, valid := security.ValidateSecureToken(cookie.Value, r.UserAgent())
    if !valid {
        http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
        return
    }

    // Récupère les infos du formulaire
    postIDStr := r.FormValue("post_id")
    content := r.FormValue("content")
    postID, err := strconv.Atoi(postIDStr)
    if err != nil || content == "" {
        http.Redirect(w, r, "/?error=badinput", http.StatusSeeOther)
        return
    }

    log.Printf("🟢 Ajouter commentaire: userID=%s, postID=%d, content=%q\n", userID, postID, content)

    // Ajoute le commentaire en base
    err = models.InsertComment(userID, postID, content)
    if err != nil {
        log.Println("❌ Erreur insertion commentaire :", err)
        http.Error(w, "Erreur lors de l'ajout du commentaire", http.StatusInternalServerError)
        return
    }

    // Redirige vers la home pour voir le résultat dans home.html
    http.Redirect(w, r, "/entry/home", http.StatusSeeOther)
}
