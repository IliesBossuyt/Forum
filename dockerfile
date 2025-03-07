# Utilise l'image officielle MySQL
FROM mysql:8.0

# Variables d'environnement pour init
ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_USER=forumuser
ENV MYSQL_PASSWORD=forumpassword
ENV MYSQL_DATABASE=forumdb

# Exposer le port MySQL
EXPOSE 3306
