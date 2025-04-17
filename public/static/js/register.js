// Fonction pour gérer l'inscription d'un utilisateur
function registerUser(event) {
    event.preventDefault();
    
    // Récupération des données du formulaire
    let data = {
        username: document.getElementById("username").value,
        email: document.getElementById("email").value,
        password: document.getElementById("password").value
    };

    // Envoi de la requête d'inscription
    fetch("/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
    .then(response => response.text())
    .then(text => {
        // Affichage du message de réponse
        const messageElement = document.getElementById("message");
        messageElement.innerText = text;
        
        // Redirection vers la page de connexion si l'inscription est réussie
        if (text.includes("réussie")) {
            messageElement.className = "success";
            setTimeout(() => window.location.href = "/auth/login", 2000);
        }
    })
    .catch(error => {
        // Gestion des erreurs réseau
        console.error("Erreur:", error);
        document.getElementById("message").innerText = "Erreur lors de l'inscription.";
    });
}