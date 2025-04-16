package models

import (
	"Forum/internal/database"
)

type Category struct {
	ID   int
	Name string
}

// Récupère toutes les catégories existantes
func GetAllCategories() ([]Category, error) {
	rows, err := database.DB.Query("SELECT id, name FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
