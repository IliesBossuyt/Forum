// Fonction de gestion des notifications
function showNotification(message, type = 'success', duration = 5000) {
    const container = document.getElementById('notification-container');
    const notification = document.createElement('div');
    const id = 'notification-' + Date.now();
    notification.id = id;
    notification.className = `notification ${type}`;

    notification.innerHTML = `
        <div class="notification-content">${message}</div>
        <button class="notification-close" onclick="removeNotification('${id}')">&times;</button>
        <div class="notification-progress" style="animation: progressBar ${duration}ms linear forwards;"></div>
    `;

    container.appendChild(notification);

    // Supprimer automatiquement après la durée spécifiée
    setTimeout(() => {
        removeNotification(id);
    }, duration);
}

function removeNotification(id) {
    const notification = document.getElementById(id);
    if (!notification) return;
    
    notification.style.animation = 'slideOut 0.3s ease-out forwards';
    setTimeout(() => {
        notification.remove();
    }, 300);
}

// Fonction de confirmation personnalisée
function showConfirm(message, options = {}) {
    return new Promise((resolve) => {
        const overlay = document.getElementById('confirm-overlay');
        const dialog = document.getElementById('confirm-dialog');
        const title = document.getElementById('confirm-title');
        const messageEl = document.getElementById('confirm-message');
        const okButton = document.getElementById('confirm-ok');
        const cancelButton = document.getElementById('confirm-cancel');
        
        // Définir le type (danger, success, warning)
        dialog.className = 'confirm-dialog';
        okButton.className = 'confirm-ok';
        
        if (options.type) {
            dialog.classList.add(options.type);
            okButton.classList.add(options.type);
        } else {
            dialog.classList.add('danger'); // Par défaut
            okButton.classList.add('danger'); 
        }
        
        // Définir titre et message
        title.textContent = options.title || 'Confirmation';
        messageEl.textContent = message;
        
        // Définir texte des boutons
        okButton.textContent = options.okText || 'Confirmer';
        cancelButton.textContent = options.cancelText || 'Annuler';
        
        // Afficher la boîte de dialogue
        overlay.style.display = 'flex';
        setTimeout(() => {
            overlay.classList.add('confirm-visible');
        }, 10);
        
        // Gestionnaires d'événements
        function handleOk() {
            overlay.classList.remove('confirm-visible');
            setTimeout(() => {
                overlay.style.display = 'none';
            }, 300);
            okButton.removeEventListener('click', handleOk);
            cancelButton.removeEventListener('click', handleCancel);
            resolve(true);
        }
        
        function handleCancel() {
            overlay.classList.remove('confirm-visible');
            setTimeout(() => {
                overlay.style.display = 'none';
            }, 300);
            okButton.removeEventListener('click', handleOk);
            cancelButton.removeEventListener('click', handleCancel);
            resolve(false);
        }
        
        okButton.addEventListener('click', handleOk);
        cancelButton.addEventListener('click', handleCancel);
    });
}

// Remplacer les alertes par notre système de notifications
document.addEventListener("DOMContentLoaded", function () {
    document.querySelectorAll(".update-role-btn").forEach(button => {
        button.addEventListener("click", function () {
            let userID = this.getAttribute("data-userid");
            let newRole = document.getElementById(`role-${userID}`).value;

            fetch('/admin/secure/change-role', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ user_id: userID, role: newRole })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showNotification("Rôle mis à jour avec succès !", "success");
                } else {
                    showNotification("Erreur : " + data.message, "error");
                }
            })
            .catch(error => {
                console.error("Erreur :", error);
                showNotification("Une erreur est survenue.", "error");
            });
        });
    });
});


// Bannir / Débannir
document.querySelectorAll(".ban-toggle-btn").forEach(button => {
    button.addEventListener("click", function () {
        let userID = this.getAttribute("data-userid");
        let isBanned = this.getAttribute("data-banned") === "true";
        let newBannedState = !isBanned;

        fetch('/admin/secure/toggle-ban', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ user_id: userID, banned: newBannedState })
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                this.innerText = newBannedState ? "Débannir" : "Bannir";
                this.setAttribute("data-banned", newBannedState.toString());
                showNotification("Utilisateur " + (newBannedState ? "banni" : "débanni") + " avec succès !", "success");
            } else {
                showNotification("Erreur : " + data.message, "error");
            }
        })
        .catch(error => {
            console.error("Erreur :", error);
            showNotification("Une erreur est survenue.", "error");
        });
    });
});

// Fonctions modales d'image inchangées
function openModal(imageSrc) {
    let modal = document.getElementById("imageModal");
    let modalImg = document.getElementById("modalContent");

    modalImg.src = imageSrc;
    modal.style.display = "flex";
    setTimeout(() => {
        modal.style.opacity = "1";
        modalImg.style.transform = "scale(1)";
    }, 10);
}

function closeModal() {
    let modal = document.getElementById("imageModal");
    let modalImg = document.getElementById("modalContent");

    modal.style.opacity = "0";
    modalImg.style.transform = "scale(0.5)";
    setTimeout(() => {
        modal.style.display = "none";
    }, 300);
}

// Supprimer un signalement
document.querySelectorAll(".delete-report-btn").forEach(button => {
    button.addEventListener("click", async () => {
        const reportIDString = button.getAttribute("data-reportid");
        const reportID = parseInt(reportIDString, 10);

        const confirmed = await showConfirm("Voulez-vous vraiment supprimer ce signalement ?", {
            title: "Supprimer le signalement",
            type: "danger",
            okText: "Supprimer"
        });
        
        if (!confirmed) return;

        fetch("/admin/delete-report-post", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ report_id: reportID }) 
        })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                showNotification("Signalement supprimé !", "success");
                location.reload();
            } else {
                showNotification("Erreur : " + data.message, "error");
            }
        });
    });
});

// Bannir l'auteur du post signalé
document.querySelectorAll(".ban-post-author-btn").forEach(button => {
    button.addEventListener("click", async () => {
        const userID = button.getAttribute("data-userid");
        
        const confirmed = await showConfirm("Voulez-vous vraiment bannir cet utilisateur ?", {
            title: "Bannir l'utilisateur",
            type: "danger",
            okText: "Bannir"
        });
        
        if (!confirmed) return;

        fetch("/admin/secure/toggle-ban", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ user_id: userID, banned: true })
        })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                showNotification("Utilisateur banni !", "success");
                location.reload();
            } else {
                showNotification("Erreur : " + data.message, "error");
            }
        });
    });
});

// Supprimer un post signalé
document.querySelectorAll(".delete-post-btn").forEach(button => {
    button.addEventListener("click", async () => {
        const postID = parseInt(button.getAttribute("data-postid"), 10);
        
        const confirmed = await showConfirm("Voulez-vous vraiment supprimer ce post, ses commentaires et signalements ?", {
            title: "Supprimer le post",
            type: "danger",
            okText: "Supprimer"
        });
        
        if (!confirmed) return;

        fetch("/admin/delete-post", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ post_id: postID.toString() })
        })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                showNotification("Post supprimé !", "success");
                location.reload();
            } else {
                showNotification("Erreur : " + data.message, "error");
            }
        })
        .catch(() => showNotification("Erreur lors de la suppression du post.", "error"));
    });
});


let currentWarnUserID = null;

function closeWarnModal() {
    document.getElementById("warnModal").style.display = "none";
    document.getElementById("warnList").innerHTML = "";
    document.getElementById("warnReason").value = "";
}

// Ouvre le casier
const userRole = document.body.dataset.userRole;
const isAdmin = userRole === "admin";

document.querySelectorAll(".warn-user-btn").forEach(button => {
    button.addEventListener("click", () => {
        const userID = button.getAttribute("data-userid");
        const username = button.getAttribute("data-username");
        currentWarnUserID = userID;

        document.getElementById("warnModalTitle").innerText = `Casier de ${username}`;

        fetch(`/admin/warns?user_id=${userID}`)
            .then(res => res.json())
            .then(data => {
                const warnList = document.getElementById("warnList");
                warnList.innerHTML = "";
                data.forEach(warn => {
                    const li = document.createElement("li");
                    const date = new Date(warn.CreatedAt).toLocaleString("fr-FR", {
                        day: "2-digit",
                        month: "2-digit",
                        year: "numeric",
                        hour: "2-digit",
                        minute: "2-digit"
                    });

                    li.innerHTML = `
                        <div><strong>Modérateur :</strong> <span>${warn.Issuer}</span></div>
                        <div><strong>Raison :</strong> <span>${warn.Reason}</span></div>
                        <div><strong>Date :</strong> <span>${date}</span></div>
                    `;

                    if (isAdmin) {
                        const deleteBtn = document.createElement("button");
                        deleteBtn.innerText = "Supprimer";
                        deleteBtn.addEventListener("click", async () => {
                            const confirmed = await showConfirm("Voulez-vous vraiment supprimer cet avertissement ?", {
                                title: "Supprimer l'avertissement",
                                type: "danger",
                                okText: "Supprimer"
                            });
                            
                            if (!confirmed) return;

                            fetch(`/admin/secure/delete-warn`, {
                                method: "POST",
                                headers: { "Content-Type": "application/json" },
                                body: JSON.stringify({ warn_id: warn.ID })
                            })
                            .then(res => res.json())
                            .then(resp => {
                                if (resp.success) {
                                    showNotification("Warn supprimé !", "success");
                                    // Supprimer le warn du DOM
                                    li.remove();

                                    // Actualiser le compteur sans recharger tout
                                    const counterSpan = document.getElementById(`warn-count-${currentWarnUserID}`);
                                    let currentCount = parseInt(counterSpan.textContent.replace(/[^\d]/g, ""));
                                    if (!isNaN(currentCount) && currentCount > 0) {
                                        counterSpan.textContent = `(${currentCount - 1})`;
                                    }
                                } else {
                                    showNotification("Erreur : " + resp.message, "error");
                                }
                            });
                        });
                        li.appendChild(deleteBtn);
                    }

                    warnList.appendChild(li);
                });

                document.getElementById("warnModal").style.display = "block";
            });
    });
});


document.getElementById("addWarnBtn").addEventListener("click", () => {
    const reason = document.getElementById("warnReason").value.trim();
    if (!reason) {
        showNotification("Raison obligatoire !", "warning");
        return;
    }

    fetch("/admin/add-warn", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ user_id: currentWarnUserID, reason: reason })
    })
    .then(res => res.json())
    .then(data => {
        if (data.success) {
            showNotification("Avertissement ajouté !", "success");
            document.getElementById(`warn-count-${currentWarnUserID}`).textContent = `(${data.new_count})`;
            closeWarnModal();
        } else {
            showNotification("Erreur : " + data.message, "error");
        }
    });
});


// Supprimer un signalement de commentaire
document.querySelectorAll(".delete-comment-report-btn").forEach(button => {
    button.addEventListener("click", async () => {
        const reportID = parseInt(button.getAttribute("data-reportid"), 10);
        
        const confirmed = await showConfirm("Voulez-vous vraiment supprimer ce signalement de commentaire ?", {
            title: "Supprimer le signalement",
            type: "danger",
            okText: "Supprimer"
        });
        
        if (!confirmed) return;

        fetch("/admin/delete-report-comment", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ report_id: reportID })
        })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                showNotification("Signalement supprimé !", "success");
                location.reload();
            } else {
                showNotification("Erreur : " + data.message, "error");
            }
        });
    });
});


// Bannir l'auteur du commentaire signalé
document.querySelectorAll(".ban-comment-author-btn").forEach(button => {
    button.addEventListener("click", async () => {
        const userID = button.getAttribute("data-userid");
        
        const confirmed = await showConfirm("Voulez-vous vraiment bannir cet utilisateur ?", {
            title: "Bannir l'utilisateur",
            type: "danger",
            okText: "Bannir"
        });
        
        if (!confirmed) return;

        fetch("/admin/secure/toggle-ban", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ user_id: userID, banned: true })
        })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                showNotification("Utilisateur banni !", "success");
                location.reload();
            } else {
                showNotification("Erreur : " + data.message, "error");
            }
        });
    });
});

// Supprimer un commentaire signalé
document.querySelectorAll(".delete-comment-btn").forEach(button => {
    button.addEventListener("click", async () => {
        const commentID = parseInt(button.getAttribute("data-commentid"), 10);
        
        const confirmed = await showConfirm("Voulez-vous vraiment supprimer ce commentaire ?", {
            title: "Supprimer le commentaire",
            type: "danger",
            okText: "Supprimer"
        });
        
        if (!confirmed) return;

        fetch("/admin/delete-comment", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ comment_id: commentID })
        })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                showNotification("Commentaire supprimé !", "success");
                location.reload();
            } else {
                showNotification("Erreur : " + data.message, "error");
            }
        })
        .catch(() => showNotification("Erreur lors de la suppression du commentaire.", "error"));
    });
});