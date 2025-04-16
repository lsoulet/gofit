package models

import "time"

type DailyMenu struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uint
	User   User      `gorm:"foreignKey:UserID"`
	Date   time.Time
	Meals  []Meal    `gorm:"many2many:dailymenu_meals;"`
}

func (d *DailyMenu) GetDailyMacroSummary() (float64, float64, float64, float64, error) {
	var totalCalories, totalProteins, totalCarbohydrates, totalLipids float64
	for _, meal := range d.Meals {
		cal, prot, carb, lipid := meal.GetMacros()
		totalCalories += cal
		totalProteins += prot
		totalCarbohydrates += carb
		totalLipids += lipid
	}
	return totalCalories, totalProteins, totalCarbohydrates, totalLipids, nil
}
