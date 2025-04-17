document.getElementById("profile-form").addEventListener("submit", function(event) {
    event.preventDefault();
    
    let data = {
        username: document.getElementById("username").value.trim(),
        email: document.getElementById("email").value.trim(),
        old_password: document.getElementById("old-password").value,
        new_password: document.getElementById("new-password").value,
        is_public: document.getElementById("is-public").checked 
    };

    const currentUsername = document.getElementById("display-username").innerText;
    fetch("/user/profile/" + currentUsername, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
    .then(response => response.text().then(text => ({ status: response.status, body: text })))
    .then(res => {
        if (res.status === 200) {
            document.getElementById("success-message").style.display = "block";
            document.getElementById("error-message").style.display = "none";

            // Mettre à jour l'affichage du username et email
            document.getElementById("display-username").innerText = data.username;
            document.getElementById("display-email").innerText = data.email;

            // Réinitialiser les champs de mot de passe
            document.getElementById("old-password").value = "";
            document.getElementById("new-password").value = "";

            // Masquer le message de succès après 3 secondes
            setTimeout(() => {
                document.getElementById("success-message").style.display = "none";
            }, 3000);
        } else {
            document.getElementById("error-message").innerText = "❌ " + res.body;
            document.getElementById("error-message").style.display = "block";
        }
    })
    .catch(error => {
        console.error("Erreur :", error);
        document.getElementById("error-message").innerText = "❌ Une erreur est survenue.";
        document.getElementById("error-message").style.display = "block";
    });
});

// Ouverture de la modale du casier
document.getElementById("openWarnBtn").addEventListener("click", function () {
    document.getElementById("warnModal").style.display = "block";
});

// Fermeture de la modale (déjà appelée dans ton HTML via `onclick`)
function closeWarnModal() {
    document.getElementById("warnModal").style.display = "none";
}