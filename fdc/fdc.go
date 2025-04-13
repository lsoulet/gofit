package fdc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	searchURL = "https://api.nal.usda.gov/fdc/v1/foods/search"
	detailURL = "https://api.nal.usda.gov/fdc/v1/food/"
)

var apiKey = os.Getenv("FDC_API_KEY") // Stocke ta clé dans une variable d'env pour la sécurité

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Foods []struct {
		Description string `json:"description"`
		FdcID       int    `json:"fdcId"`
	} `json:"foods"`
}

type FoodNutrient struct {
	Name  string  `json:"name"`
	Unit  string  `json:"unitName"`
	Value float64 `json:"value"`
}

type FoodDetail struct {
	Description   string         `json:"description"`
	FoodNutrients []FoodNutrient `json:"foodNutrients"`
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
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var results []string
	for _, food := range result.Foods {
		results = append(results, fmt.Sprintf("%s (fdcId: %d)", food.Description, food.FdcID))
	}
	return results, nil
}

// 📋 Détails nutritionnels
func GetFoodDetails(fdcID int) (map[string]float64, error) {
	url := fmt.Sprintf("%s%d?api_key=%s", detailURL, fdcID, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var detail FoodDetail
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, err
	}

	nutrients := map[string]float64{}
	for _, nutrient := range detail.FoodNutrients {
		switch nutrient.Name {
		case "Energy":
			nutrients["calories"] = nutrient.Value
		case "Protein":
			nutrients["protein"] = nutrient.Value
		case "Total lipid (fat)":
			nutrients["fat"] = nutrient.Value
		case "Carbohydrate, by difference":
			nutrients["carbohydrates"] = nutrient.Value
		case "Fiber, total dietary":
			nutrients["fiber"] = nutrient.Value
		}
	}
	return nutrients, nil
}
