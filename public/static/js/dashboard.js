// Fonction pour afficher une notification
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

    // Supprime automatiquement la notification après la durée spécifiée
    setTimeout(() => {
        removeNotification(id);
    }, duration);
}

// Fonction pour supprimer une notification
function removeNotification(id) {
    const notification = document.getElementById(id);
    if (!notification) return;
    
    notification.style.animation = 'slideOut 0.3s ease-out forwards';
    setTimeout(() => {
        notification.remove();
    }, 300);
}

// Fonction pour afficher une boîte de dialogue de confirmation
function showConfirm(message, options = {}) {
    return new Promise((resolve) => {
        const overlay = document.getElementById('confirm-overlay');
        const dialog = document.getElementById('confirm-dialog');
        const title = document.getElementById('confirm-title');
        const messageEl = document.getElementById('confirm-message');
        const okButton = document.getElementById('confirm-ok');
        const cancelButton = document.getElementById('confirm-cancel');
        
        // Configure le style de la boîte de dialogue
        dialog.className = 'confirm-dialog';
        okButton.className = 'confirm-ok';
        
        if (options.type) {
            dialog.classList.add(options.type);
            okButton.classList.add(options.type);
        } else {
            dialog.classList.add('danger'); // Style par défaut
            okButton.classList.add('danger'); 
        }
        
        // Configure le contenu de la boîte de dialogue
        title.textContent = options.title || 'Confirmation';
        messageEl.textContent = message;
        okButton.textContent = options.okText || 'Confirmer';
        cancelButton.textContent = options.cancelText || 'Annuler';
        
        // Affiche la boîte de dialogue
        overlay.style.display = 'flex';
        setTimeout(() => {
            overlay.classList.add('confirm-visible');
        }, 10);
        
        // Gestionnaires d'événements pour les boutons
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

// Gestion des mises à jour de rôle
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

// Gestion du bannissement/débannissement
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

// Fonctions pour la modale d'image
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

// Gestion de la suppression des signalements
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

// Gestion du bannissement des auteurs de posts signalés
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

// Gestion de la suppression des posts signalés
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

// Gestion des avertissements
let currentWarnUserID = null;

function closeWarnModal() {
    document.getElementById("warnModal").style.display = "none";
    document.getElementById("warnList").innerHTML = "";
    document.getElementById("warnReason").value = "";
}

// Vérifie si l'utilisateur est administrateur
const userRole = document.body.dataset.userRole;
const isAdmin = userRole === "admin";

// Gestion de l'affichage du casier d'avertissements
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
                                    li.remove();
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

// Gestion de l'ajout d'avertissements
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

// Gestion de la suppression des signalements de commentaires
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

// Gestion du bannissement des auteurs de commentaires signalés
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

// Gestion de la suppression des commentaires signalés
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