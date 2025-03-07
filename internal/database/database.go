package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Fonction pour initialiser la base de données
func InitDatabase() {
	var err error
	DB, err = sql.Open("sqlite3", "../forum.db")
	if err != nil {
		log.Fatal("Erreur de connexion à SQLite :", err)
	}

	// Vérifier la connexion
	err = DB.Ping()
	if err != nil {
		log.Fatal("Impossible de se connecter à la base de données :", err)
	}

	log.Println("Connexion SQLite réussie !")
}
