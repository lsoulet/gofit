package fdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/lsoulet/gofit/db"
	"github.com/lsoulet/gofit/models"
)

const (
	searchURL = "https://api.nal.usda.gov/fdc/v1/foods/search"
	detailURL = "https://api.nal.usda.gov/fdc/v1/food/"
)

var apiKey = os.Getenv("FDC_API_KEY") // Utiliser une variable d'env pour plus de sécurité

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Foods []struct {
		Description string `json:"description"`
		FdcID       int    `json:"fdcId"`
	} `json:"foods"`
}
type FoodDetail struct {
	Description   string `json:"description"`
	FoodNutrients []struct {
		Nutrient struct {
			Number string `json:"number"`
			Name   string `json:"name"`
		} `json:"nutrient"`
		Amount float64 `json:"amount"`
	} `json:"foodNutrients"`
}

// 🔍 Rechercher un aliment
func SearchFood(query string) ([]string, error) {
	reqBody, _ := json.Marshal(SearchRequest{Query: query})

	req, err := http.NewRequest("POST", fmt.Sprintf("%s?api_key=%s", searchURL, apiKey), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var results []string
	for _, food := range result.Foods {
		results = append(results, fmt.Sprintf("%s (fdcId: %d)", food.Description, food.FdcID))
	}
	return results, nil
}

// 🔎 Récupérer les détails nutritionnels d'un aliment à partir de son fdcId
func GetFoodDetails(fdcID int) (string, float64, float64, float64, float64, error) {
	url := fmt.Sprintf("%s%d?api_key=%s", detailURL, fdcID, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return "", 0, 0, 0, 0, err
	}
	defer resp.Body.Close()

	var result FoodDetail
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", 0, 0, 0, 0, err
	}

	var calories, proteins, carbohydrates, lipids float64

	for _, nutrient := range result.FoodNutrients {
		switch nutrient.Nutrient.Number {
		case "203": // Protéines
			proteins = nutrient.Amount
		case "204": // Lipides
			lipids = nutrient.Amount
		case "205": // Glucides
			carbohydrates = nutrient.Amount
		case "208": // Calories
			calories = nutrient.Amount
		}
	}

	return result.Description, calories, proteins, carbohydrates, lipids, nil
}

func AddFoodToMeal(mealID uint, fdcID int, quantity float64) error {
	// Récupérer le repas
	var meal models.Meal
	if err := db.DB.First(&meal, mealID).Error; err != nil {
		return fmt.Errorf("erreur lors de la récupération du repas : %w", err)
	}

	// Récupérer les détails de l'aliment
	name, calories, proteins, carbs, lipids, err := GetFoodDetails(fdcID)
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération des détails de l'aliment : %w", err)
	}

	// Calculer les valeurs nutritionnelles en fonction de la quantité
	ratio := quantity / 100.0
	meal.Calories += calories * ratio
	meal.Proteins += proteins * ratio
	meal.Carbohydrates += carbs * ratio
	meal.Lipids += lipids * ratio

	// Sauvegarder les modifications
	if err := db.DB.Save(&meal).Error; err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du repas : %w", err)
	}

	fmt.Printf("✔ Aliment '%s' (%.0f g) ajouté au repas\n", name, quantity)
	return nil
}

func GetMeals() ([]models.Meal, error) {
	var meals []models.Meal
	return meals, nil
}

func ListMeals() error {
	meals, err := GetMeals()
	if err != nil {
		return err
	}

	if len(meals) == 0 {
		fmt.Println("Aucun repas enregistré.")
		return nil
	}

	fmt.Println("🍽️ Repas enregistrés :")
	for i, meal := range meals {
		fmt.Printf("%d. %s (%s) | %.1f kcal | P: %.1f g | G: %.1f g | L: %.1f g\n",
			i+1, meal.Description, meal.Type, meal.Calories, meal.Proteins, meal.Carbohydrates, meal.Lipids)
	}
	return nil
}

func AddMeal(meal models.Meal) error {
	return db.DB.Create(&meal).Error
}
