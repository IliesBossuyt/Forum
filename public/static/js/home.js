// Ouvrir le formulaire de création de post
function openCreatePostForm() {
    document.getElementById("create-post-form").style.display = "flex";
    setTimeout(() => {
        document.getElementById("create-post-form").style.opacity = "1";
        document.querySelector("#create-post-form .modal-content").style.transform = "scale(1)";
    }, 10);
}

function closeCreatePostForm() {
    document.getElementById("create-post-form").style.opacity = "0";
    document.querySelector("#create-post-form .modal-content").style.transform = "scale(0.8)";
    setTimeout(() => {
        document.getElementById("create-post-form").style.display = "none";
    }, 300);
}

// Like / Dislike d'un post
function likePost(postID, value) {
    fetch('/user/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ post_id: parseInt(postID), value: value })
    })
    .then(response => response.json())
    .then(data => {
        document.getElementById(`likes-${postID}`).innerText = data.likes || 0;
        document.getElementById(`dislikes-${postID}`).innerText = data.dislikes || 0;
    })
    .catch(error => {
        console.error('Erreur:', error);
    });
}

// Ouvrir le formulaire de suppression
function openDeleteForm(postID) {
    document.getElementById("delete-post-id").value = postID;
    document.getElementById("delete-form").style.display = "flex";
    setTimeout(() => {
        document.getElementById("delete-form").style.opacity = "1";
        document.querySelector("#delete-form .modal-content").style.transform = "scale(1)";
    }, 10);
}

function closeDeleteForm() {
    document.getElementById("delete-form").style.opacity = "0";
    document.querySelector("#delete-form .modal-content").style.transform = "scale(0.8)";
    setTimeout(() => {
        document.getElementById("delete-form").style.display = "none";
    }, 300);
}

// Confirme la suppression du post
function confirmDelete() {
    let postID = document.getElementById("delete-post-id").value;
    deletePost(postID);
    closeDeleteForm();
}

// Supprime le post
function deletePost(postID) {
    fetch('/user/delete-post', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ post_id: postID })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            let postElement = document.getElementById(`post-wrapper-${postID}`);
            if (postElement) {
                postElement.remove();
            }
            showNotification("Post supprimé avec succès", "success");
        } else {
            showNotification(data.message || "Erreur lors de la suppression du post", "error");
        }
    })
    .catch(error => {
        console.error("Erreur JS :", error);
        showNotification("Erreur réseau lors de la suppression du post", "error");
    });
}

function editPost(postID) {
    let currentContent = document.getElementById(`content-${postID}`).innerText;
    document.getElementById("edit-content").value = currentContent;
    document.getElementById("edit-post-id").value = postID;
    document.getElementById("delete-image").checked = false;
    document.getElementById("edit-form").style.display = "flex";
    setTimeout(() => {
        document.getElementById("edit-form").style.opacity = "1";
        document.querySelector("#edit-form .modal-content").style.transform = "scale(1)";
    }, 10);
}

function closeEditForm() {
    document.getElementById("edit-form").style.opacity = "0";
    document.querySelector("#edit-form .modal-content").style.transform = "scale(0.8)";
    setTimeout(() => {
    document.getElementById("edit-form").style.display = "none";
    }, 300);
}

function submitEdit() {
    let postID = document.getElementById("edit-post-id").value;
    let newContent = document.getElementById("edit-content").value;
    let deleteImage = document.getElementById("delete-image").checked;
    let formData = new FormData();
    formData.append("post_id", postID);
    formData.append("content", newContent);
    formData.append("delete_image", deleteImage);

    let fileInput = document.getElementById("edit-image");
    if (fileInput.files.length > 0) {
        formData.append("image", fileInput.files[0]);
    }

    fetch('/user/edit-post', {
        method: 'POST',
        body: formData
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            // Mise à jour du texte du post
            document.getElementById(`content-${postID}`).innerText = newContent;

            // Gestion de l'image
            if (data.imageUpdated) {
                let imgContainer = document.querySelector(`#post-${postID} .post-image`);
                let imgElement = imgContainer ? imgContainer.querySelector('img') : null;
                
            if (!imgElement) {
                    // Créer le conteneur d'image s'il n'existe pas
                    if (!imgContainer) {
                        imgContainer = document.createElement('div');
                        imgContainer.className = 'post-image';
                        let contentElement = document.getElementById(`content-${postID}`);
                        contentElement.parentNode.insertBefore(imgContainer, contentElement.nextSibling);
                    }
                    
                    // Créer l'élément image
                    imgElement = document.createElement('img');
                imgElement.alt = "Image postée";
                imgElement.onclick = function() { openModal(`/entry/image/${postID}`); };
                    imgContainer.appendChild(imgElement);
                }
                
                // Mettre à jour l'image avec un timestamp unique
            imgElement.src = `/entry/image/${postID}?t=${new Date().getTime()}`;
            } 
            
            if (data.imageDeleted) {
                let imgContainer = document.querySelector(`#post-${postID} .post-image`);
                if (imgContainer) {
                    imgContainer.remove();
                }
            }

            closeEditForm();
            showNotification("Post modifié avec succès", "success");
        } else {
            showNotification(data.message || "Erreur lors de la modification du post", "error");
        }
    })
    .catch(error => {
        console.error("Erreur lors de la modification :", error);
        showNotification("Erreur réseau lors de la modification du post", "error");
    });
}

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

// Fonction modifiée pour ouvrir le modal de signalement au lieu du prompt
function reportPost(postID) {
    document.getElementById("report-post-id").value = postID;
    document.getElementById("report-reason").value = "";
    document.getElementById("report-post-form").style.display = "flex";
    setTimeout(() => {
        document.getElementById("report-post-form").style.opacity = "1";
        document.querySelector("#report-post-form .modal-content").style.transform = "scale(1)";
    }, 10);
}

function closeReportForm() {
    document.getElementById("report-post-form").style.opacity = "0";
    document.querySelector("#report-post-form .modal-content").style.transform = "scale(0.8)";
    setTimeout(() => {
        document.getElementById("report-post-form").style.display = "none";
    }, 300);
}

function submitReport() {
    const postID = document.getElementById("report-post-id").value;
    const reason = document.getElementById("report-reason").value.trim();

    const payload = {
        post_id: Number(postID),
        reason: reason
    };

    fetch("/user/report", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
    })
    .then(res => res.json())
    .then(data => {
        if (data.success) {
            closeReportForm();
            showNotification("Le post a été signalé avec succès", "success");
        } else {
            showNotification("Une erreur est survenue lors du signalement", "error");
        }
    })
    .catch(() => showNotification("Erreur réseau", "error"));
}

function toggleComments(postID) {
    const section = document.getElementById("comments-" + postID);
    if (section.style.display === "none") {
        section.style.display = "block";
    } else {
        section.style.display = "none";
    }
}

function submitComment(postID) {
    const textarea = document.getElementById("comment-content-" + postID);
    const content = textarea.value.trim();
    if (!content) {
        showNotification("Veuillez écrire un commentaire", "warning");
        return;
    }

    fetch("/user/comment", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            post_id: parseInt(postID),
            content: content
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            // Recharger la page pour voir le nouveau commentaire
            setTimeout(() => {
                location.reload();
            }, 1000);
    } else {
            showNotification(data.message || "Erreur lors de l'ajout du commentaire", "error");
    }
})
    .catch(error => {
        console.error("Erreur:", error);
        showNotification("Erreur réseau lors de l'ajout du commentaire", "error");
    });
}

function likeComment(commentID, value) {
    fetch('/user/like-comment', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ comment_id: parseInt(commentID), value: value })
    })
    .then(response => response.json())
    .then(data => {
        document.getElementById(`comment-like-${commentID}`).innerText = data.likes || 0;
        document.getElementById(`comment-dislike-${commentID}`).innerText = data.dislikes || 0;
    })
    .catch(error => {
        console.error('Erreur:', error);
    });
}

function reportComment(commentID) {
    document.getElementById("report-comment-id").value = commentID;
    document.getElementById("report-comment-reason").value = "";
    document.getElementById("report-comment-form").style.display = "flex";
    setTimeout(() => {
        document.getElementById("report-comment-form").style.opacity = "1";
        document.querySelector("#report-comment-form .modal-content").style.transform = "scale(1)";
    }, 10);
}

function closeReportCommentForm() {
    document.getElementById("report-comment-form").style.opacity = "0";
    document.querySelector("#report-comment-form .modal-content").style.transform = "scale(0.8)";
    setTimeout(() => {
        document.getElementById("report-comment-form").style.display = "none";
    }, 300);
}

function submitCommentReport() {
    const commentID = document.getElementById("report-comment-id").value;
    const reason = document.getElementById("report-comment-reason").value.trim();
    
    fetch("/user/report-comment", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ comment_id: Number(commentID), reason: reason })
    })
    .then(res => res.json())
    .then(data => {
        if (data.success) {
            closeReportCommentForm();
            showNotification("Le commentaire a été signalé avec succès", "success");
        } else {
            showNotification("Une erreur est survenue lors du signalement", "error");
        }
    })
    .catch(() => showNotification("Erreur réseau", "error"));
}

function editComment(commentID, content) {
    document.getElementById("edit-comment-content").value = content;
    document.getElementById("edit-comment-id").value = commentID;
    document.getElementById("edit-comment-form").style.display = "flex";
    setTimeout(() => {
        document.getElementById("edit-comment-form").style.opacity = "1";
        document.querySelector("#edit-comment-form .modal-content").style.transform = "scale(1)";
    }, 10);
}

function closeEditCommentForm() {
    document.getElementById("edit-comment-form").style.opacity = "0";
    document.querySelector("#edit-comment-form .modal-content").style.transform = "scale(0.8)";
    setTimeout(() => {
    document.getElementById("edit-comment-form").style.display = "none";
    }, 300);
}

function submitEditComment() {
    const commentID = document.getElementById("edit-comment-id").value;
    const content = document.getElementById("edit-comment-content").value.trim();

    if (!content) {
        showNotification("Le commentaire ne peut pas être vide", "warning");
        return;
    }

    fetch("/user/edit-comment", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ comment_id: Number(commentID), content: content })
    })
    .then(res => res.json())
    .then(data => {
        if (data.success) {
            // Mise à jour du contenu du commentaire dans le DOM
            const commentContentElement = document.querySelector(`#comment-${commentID} .comment-content`);
            if (commentContentElement) {
                commentContentElement.textContent = content;
            }
            closeEditCommentForm();
            showNotification("Commentaire modifié avec succès", "success");
        } else {
            showNotification(data.message || "Erreur lors de la modification du commentaire", "error");
        }
    })
    .catch(() => {
        showNotification("Erreur réseau lors de la modification du commentaire", "error");
    });
}

function deleteComment(commentID) {
    document.getElementById("delete-comment-id").value = commentID;
    document.getElementById("delete-comment-form").style.display = "flex";
    setTimeout(() => {
        document.getElementById("delete-comment-form").style.opacity = "1";
        document.querySelector("#delete-comment-form .modal-content").style.transform = "scale(1)";
    }, 10);
}

function closeDeleteCommentForm() {
    document.getElementById("delete-comment-form").style.opacity = "0";
    document.querySelector("#delete-comment-form .modal-content").style.transform = "scale(0.8)";
    setTimeout(() => {
        document.getElementById("delete-comment-form").style.display = "none";
    }, 300);
}

function confirmDeleteComment() {
    const commentID = document.getElementById("delete-comment-id").value;
    
    fetch("/user/delete-comment", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ comment_id: parseInt(commentID) })
    })
    .then(res => res.json())
    .then(data => {
        if (data.success) {
            // Supprimer le commentaire du DOM
            const commentElement = document.getElementById(`comment-${commentID}`);
            if (commentElement) {
                commentElement.remove();
            }
            closeDeleteCommentForm();
            showNotification("Commentaire supprimé avec succès", "success");
        } else {
            showNotification(data.message || "Erreur lors de la suppression du commentaire", "error");
        }
    })
    .catch(() => {
        showNotification("Erreur réseau lors de la suppression du commentaire", "error");
    });
}

// Notifications
function toggleNotifMenu() {
    const dropdown = document.getElementById("notifDropdown");
    if (dropdown.style.display === "none") {
        fetchNotifications();
        dropdown.style.display = "block";
    } else {
        dropdown.style.display = "none";
    }
}

function fetchNotifications() {
    fetch("/user/notifications")
        .then(res => {
            if (!res.ok) {
                throw new Error(`Erreur HTTP: ${res.status}`);
            }
            return res.json();
        })
        .then(data => {
            const list = document.getElementById("notifList");
            list.innerHTML = "";

            if (!Array.isArray(data)) {
                console.error("Format de données inattendu:", data);
                list.innerHTML = "<li>Format de données inattendu</li>";
                return;
            }
            
            if (data.length === 0) {
                const li = document.createElement("li");
                li.textContent = "Aucune notification.";
                list.appendChild(li);
                document.getElementById("notifCount").style.display = "none";
                return;
            }
            
            data.forEach(notif => {
                const li = document.createElement("li");
                if (!notif.seen) {
                    li.classList.add("unseen");
                }
                
                // Utiliser directement le message fourni par le serveur
                const message = notif.message || "Nouvelle notification";
                
                // Formater la date ou utiliser une valeur par défaut
                let dateStr = "Date inconnue";
                try {
                    if (notif.created_at) {
                        dateStr = notif.created_at;
                    }
                } catch (e) {
                    console.error("Erreur de formatage de date", e);
                }
                
                li.innerHTML = `
                    <div>${message}</div>
                    <small>${dateStr}</small>
                `;
                
                list.appendChild(li);
            });
            
            // Mettre à jour le compteur de notifications non vues
            const unseenCount = data.filter(notif => !notif.seen).length;
            const countElement = document.getElementById("notifCount");
            if (unseenCount > 0) {
                countElement.textContent = unseenCount;
                countElement.style.display = "inline";
            } else {
                countElement.style.display = "none";
            }
        })
        .catch(error => {
            console.error("Erreur de chargement des notifications:", error);
            const list = document.getElementById("notifList");
            list.innerHTML = "<li>Erreur lors du chargement des notifications.</li>";
            
            // S'assurer que le compteur est masqué en cas d'erreur
            const countElement = document.getElementById("notifCount");
            if (countElement) {
                countElement.style.display = "none";
            }
        });
}

function markNotificationsRead() {
    fetch("/user/notifications/mark-read", {
        method: "POST",
        headers: { "Content-Type": "application/json" }
    })
    .then(() => {
        // Enlever la classe "unseen" de tous les éléments
        document.querySelectorAll("#notifList li").forEach(li => {
            li.classList.remove("unseen");
        });
        fetchNotifications();
    });
}

function deleteAllNotifications() {
    fetch("/user/notifications/delete-all", {
        method: "POST",
        headers: { "Content-Type": "application/json" }
    })
    .then(() => {
        // après suppression, on remplace manuellement la liste
        document.getElementById("notifList").innerHTML = "<li>Aucune notification.</li>";
        document.getElementById("notifCount").style.display = "none";
    });
}

document.addEventListener("DOMContentLoaded", () => {
    fetchNotifications();
});




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

// Pour compatibilité, mapper l'ancienne fonction showToast vers la nouvelle
function showToast(message, type = "info") {
    showNotification(message, type);
}
