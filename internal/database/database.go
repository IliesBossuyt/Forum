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
	DB, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal("Erreur de connexion à SQLite :", err)
	}

	// Vérifier la connexion
	err = DB.Ping()
	if err != nil {
		log.Fatal("Impossible de se connecter à la base de données :", err)
	}

	log.Println("Connexion SQLite réussie !")

	// Créer les tables
	createTables()
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT DEFAULT 'user'
	);
	`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Erreur lors de la création des tables :", err)
	}
}
