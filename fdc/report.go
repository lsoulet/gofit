package fdc

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/lsoulet/gofit/db"
	"github.com/lsoulet/gofit/models"
)

// GenerateNutritionalReport génère un rapport nutritionnel pour tous les menus journaliers
func GenerateNutritionalReport() error {
	var menus []models.DailyMenu
	if err := db.DB.Preload("User").Preload("Meals").Find(&menus).Error; err != nil {
		return fmt.Errorf("erreur lors de la récupération des menus : %w", err)
	}

	if len(menus) == 0 {
		fmt.Println("Aucun menu journalier enregistré.")
		return nil
	}

	// Créer le tableau
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Utilisateur", "Calories", "Protéines", "Glucides", "Lipides"})

	// Pour chaque menu, calculer les totaux
	for _, menu := range menus {
		var totalCal, totalProt, totalCarb, totalLip float64
		for _, meal := range menu.Meals {
			totalCal += meal.Calories
			totalProt += meal.Proteins
			totalCarb += meal.Carbohydrates
			totalLip += meal.Lipids
		}

		// Ajouter une ligne au tableau
		table.Append([]string{
			menu.Date.Format("02/01/2006"),
			fmt.Sprintf("%s %s", menu.User.FirstName, menu.User.LastName),
			strconv.FormatFloat(totalCal, 'f', 1, 64),
			strconv.FormatFloat(totalProt, 'f', 1, 64),
			strconv.FormatFloat(totalCarb, 'f', 1, 64),
			strconv.FormatFloat(totalLip, 'f', 1, 64),
		})
	}

	// Afficher le tableau
	fmt.Println("\n📊 Rapport nutritionnel :")
	table.Render()
	return nil
}
