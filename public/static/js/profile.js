// Gestion du formulaire de profil
document.getElementById("profile-form").addEventListener("submit", function(event) {
    event.preventDefault();
    
    // Récupération des données du formulaire
    let data = {
        username: document.getElementById("username").value.trim(),
        email: document.getElementById("email").value.trim(),
        old_password: document.getElementById("old-password").value,
        new_password: document.getElementById("new-password").value,
        is_public: document.getElementById("is-public").checked 
    };

    // Récupération du nom d'utilisateur actuel
    const currentUsername = document.getElementById("display-username").innerText;
    
    // Envoi de la requête de mise à jour du profil
    fetch("/user/profile/" + currentUsername, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
    .then(response => response.text().then(text => ({ status: response.status, body: text })))
    .then(res => {
        // Traitement de la réponse en cas de succès
        if (res.status === 200) {
            document.getElementById("success-message").style.display = "block";
            document.getElementById("error-message").style.display = "none";

            // Mise à jour des informations affichées
            document.getElementById("display-username").innerText = data.username;
            document.getElementById("display-email").innerText = data.email;

            // Réinitialisation des champs de mot de passe
            document.getElementById("old-password").value = "";
            document.getElementById("new-password").value = "";

            // Masquage automatique du message de succès
            setTimeout(() => {
                document.getElementById("success-message").style.display = "none";
            }, 3000);
        } else {
            // Affichage du message d'erreur
            document.getElementById("error-message").innerText = "❌ " + res.body;
            document.getElementById("error-message").style.display = "block";
        }
    })
    .catch(error => {
        // Gestion des erreurs réseau
        console.error("Erreur :", error);
        document.getElementById("error-message").innerText = "❌ Une erreur est survenue.";
        document.getElementById("error-message").style.display = "block";
    });
});

// Gestion de l'ouverture de la modale des avertissements
document.getElementById("openWarnBtn").addEventListener("click", function () {
    document.getElementById("warnModal").style.display = "block";
});

// Gestion de la fermeture de la modale des avertissements
function closeWarnModal() {
    document.getElementById("warnModal").style.display = "none";
}