// Fonction pour gérer la connexion d'un utilisateur
function loginUser(event) {
    event.preventDefault();

    let data = {
        identifier: document.getElementById("identifier").value,
        password: document.getElementById("password").value
    };

    // Envoie une requête POST pour la connexion
    fetch("/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
    .then(async response => {
        // Vérifie si l'utilisateur est banni
        if (response.status === 403) {
            const res = await response.json();
            if (res.banned === "true") {
                const messageElement = document.getElementById("message");
                messageElement.innerText = "⚠️ Vous avez été banni du site. Veuillez contacter un administrateur.";
                messageElement.className = "error";
                return;
            }
        }

        // Traite la réponse de la connexion
        const res = await response.json();
        const messageElement = document.getElementById("message");
        messageElement.innerText = res.message;
        
        // Redirige vers la page d'accueil si la connexion est réussie
        if (res.message.includes("réussie")) {
            messageElement.className = "success";
            setTimeout(() => window.location.href = "/entry/home", 500);
        } else {
            messageElement.className = "error";
        }
    })
    .catch(error => {
        console.error("Erreur:", error);
        const messageElement = document.getElementById("message");
        messageElement.innerText = "Erreur lors de la connexion.";
        messageElement.className = "error";
    });
}

// Fonction pour ouvrir la fenêtre de connexion Google
function openGoogleLogin() {
    const width = 600, height = 700;
    const left = (screen.width - width) / 2;
    const top = (screen.height - height) / 2;

    // Ouvre une nouvelle fenêtre pour l'authentification Google
    window.open(
        "https://localhost:8443/auth/google/login",  // L'URL d'authentification Google
        "GoogleLoginPopup",
        `width=${width},height=${height},top=${top},left=${left},resizable=yes,scrollbars=yes,status=yes`
    );
}

// Fonction pour ouvrir la fenêtre de connexion GitHub
function openGitHubLogin() {
    const width = 600, height = 700;
    const left = (screen.width - width) / 2;
    const top = (screen.height - height) / 2;

    // Ouvre une nouvelle fenêtre pour l'authentification GitHub
    window.open("/auth/github/login", "GitHub Login",
        `width=${width},height=${height},top=${top},left=${left}`);
}