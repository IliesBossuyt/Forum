function loginUser(event) {
    event.preventDefault();

    let data = {
        identifier: document.getElementById("identifier").value,
        password: document.getElementById("password").value
    };

    fetch("/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
    .then(async response => {
    if (response.status === 403) {
        const res = await response.json();
        if (res.banned === "true") {
            const messageElement = document.getElementById("message");
            messageElement.innerText = "⚠️ Vous avez été banni du site. Veuillez contacter un administrateur.";
            messageElement.className = "error";
            return;
        }
    }

    const res = await response.json();
    const messageElement = document.getElementById("message");
    messageElement.innerText = res.message;
    
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


    function openGoogleLogin() {
        const width = 600, height = 700;
        const left = (screen.width - width) / 2;
        const top = (screen.height - height) / 2;
    
        window.open(
            "https://localhost:8443/auth/google/login",  // L'URL d'authentification Google
            "GoogleLoginPopup",
            `width=${width},height=${height},top=${top},left=${left},resizable=yes,scrollbars=yes,status=yes`
        );
    }

    function openGitHubLogin() {
    const width = 600, height = 700;
    const left = (screen.width - width) / 2;
    const top = (screen.height - height) / 2;

    window.open("/auth/github/login", "GitHub Login",
        `width=${width},height=${height},top=${top},left=${left}`);
}