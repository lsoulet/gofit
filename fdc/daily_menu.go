package fdc

import (
	"fmt"
	"time"

	"github.com/lsoulet/gofit/db"
	"github.com/lsoulet/gofit/models"
)

// GetUsers récupère la liste des utilisateurs
func GetUsers() ([]models.User, error) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des utilisateurs : %w", err)
	}
	return users, nil
}

// GetDailyMenus récupère la liste des menus journaliers
func GetDailyMenus() ([]models.DailyMenu, error) {
	var menus []models.DailyMenu
	if err := db.DB.Preload("User").Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des menus : %w", err)
	}
	return menus, nil
}

// CreateDailyMenu crée un nouveau menu journalier
func CreateDailyMenu(userID uint, date time.Time) error {
	menu := models.DailyMenu{
		UserID: userID,
		Date:   date,
	}
	if err := db.DB.Create(&menu).Error; err != nil {
		return fmt.Errorf("erreur lors de la création du menu : %w", err)
	}
	return nil
}

// AddMealToDailyMenu ajoute un repas à un menu journalier
func AddMealToDailyMenu(menuID uint, mealType models.MealType, description string) error {
	// Récupérer le menu
	var menu models.DailyMenu
	if err := db.DB.First(&menu, menuID).Error; err != nil {
		return fmt.Errorf("erreur lors de la récupération du menu : %w", err)
	}

	// Récupérer le repas source
	var sourceMeal models.Meal
	if err := db.DB.Where("description = ? AND type = ?", description, mealType).First(&sourceMeal).Error; err != nil {
		return fmt.Errorf("erreur lors de la récupération du repas source : %w", err)
	}

	// Vérifier si un repas de ce type existe déjà (sauf pour les collations)
	if mealType != models.Snack {
		var existingMenu models.DailyMenu
		if err := db.DB.Preload("Meals").First(&existingMenu, menuID).Error; err != nil {
			return fmt.Errorf("erreur lors de la vérification des repas existants : %w", err)
		}

		for _, meal := range existingMenu.Meals {
			if meal.Type == mealType {
				return fmt.Errorf("ce menu contient déjà un repas de type %s", mealType)
			}
		}
	}

	// Créer le nouveau repas avec les valeurs nutritionnelles du repas source
	meal := models.Meal{
		Type:          mealType,
		Description:   description,
		Calories:      sourceMeal.Calories,
		Proteins:      sourceMeal.Proteins,
		Carbohydrates: sourceMeal.Carbohydrates,
		Lipids:        sourceMeal.Lipids,
	}

	// Sauvegarder le repas
	if err := db.DB.Create(&meal).Error; err != nil {
		return fmt.Errorf("erreur lors de la création du repas : %w", err)
	}

	// Associer le repas au menu
	if err := db.DB.Model(&menu).Association("Meals").Append(&meal); err != nil {
		return fmt.Errorf("erreur lors de l'association du repas au menu : %w", err)
	}

	return nil
}

// ListDailyMenus affiche la liste des menus journaliers
func ListDailyMenus() error {
	menus, err := GetDailyMenus()
	if err != nil {
		return err
	}

	if len(menus) == 0 {
		fmt.Println("Aucun menu journalier enregistré.")
		return nil
	}

	fmt.Println("📅 Menus journaliers enregistrés :")
	for i, menu := range menus {
		fmt.Printf("%d. %s %s - %s\n", i+1, menu.User.FirstName, menu.User.LastName, menu.Date.Format("02/01/2006"))
	}
	return nil
}
