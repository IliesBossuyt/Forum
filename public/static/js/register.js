function registerUser(event) {
    event.preventDefault();
    let data = {
        username: document.getElementById("username").value,
        email: document.getElementById("email").value,
        password: document.getElementById("password").value
    };

    fetch("/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
    .then(response => response.text())
    .then(text => {
        const messageElement = document.getElementById("message");
        messageElement.innerText = text;
        
        if (text.includes("réussie")) {
            messageElement.className = "success";
            setTimeout(() => window.location.href = "/auth/login", 2000);
        }
    })
    .catch(error => {
        console.error("Erreur:", error);
        document.getElementById("message").innerText = "Erreur lors de l'inscription.";
    });
}