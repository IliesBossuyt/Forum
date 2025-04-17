package models

import (
	"Forum/internal/database"
)

// Structure d'une catégorie
type Category struct {
	ID   int    // Identifiant unique
	Name string // Nom de la catégorie
}

// Récupère toutes les catégories
func GetAllCategories() ([]Category, error) {
	// Exécute la requête SQL
	rows, err := database.DB.Query("SELECT id, name FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parcourt les résultats
	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
