package models

import "time"

type DailyMenu struct {
	Date  time.Time // Unique per day
	Meals []Meal
}

func (d *DailyMenu) GetDailyMacroSummary() (float64, float64, float64, float64, error) {
	var totalCalories, totalProteins, totalCarbohydrates, totalLipids float64
	for _, meal := range d.Meals {
		cal, prot, carb, lipid := meal.CalculateMacros()
		totalCalories += cal
		totalProteins += prot
		totalCarbohydrates += carb
		totalLipids += lipid
	}
	return totalCalories, totalProteins, totalCarbohydrates, totalLipids, nil
}
