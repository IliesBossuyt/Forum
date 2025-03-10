package database

import (
    "database/sql"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDatabase() {
    dsn := "forumuser:forumpassword@tcp(127.0.0.1:3306)/forumdb?charset=utf8mb4&parseTime=True&loc=Local"
    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Erreur de connexion MySQL :", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Impossible de se connecter à MySQL :", err)
    }

    log.Println("Connexion MySQL réussie !")
}

