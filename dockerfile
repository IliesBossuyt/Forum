# On part de l'image Alpine
FROM alpine:latest

# Installer sqlite
RUN apk add --no-cache sqlite

# Définir un dossier de travail
WORKDIR /db

# Par défaut, on lance la CLI sqlite3
ENTRYPOINT ["sqlite3"]

# Commande par défaut (pas d'argument => pas de DB)
CMD []
