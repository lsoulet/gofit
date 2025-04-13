package models

import (
	"errors"
	"math"
	"os"
	"sort"
	"time"

	chart "github.com/wcharczuk/go-chart/v2"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type User struct {
	FirstName         string
	LastName          string
	Age               int
	Measurements      []Measurement
	Goal              Goal
	Gender            Gender
	DailyMenus        []DailyMenu
	CalorieNeeds      float64
	ProteinNeeds      float64
	CarohydratesNeeds float64
	LipidNeeds        float64
}

func (u *User) GetMealsByDate(date time.Time) ([]Meal, error) {
	for _, dm := range u.DailyMenus {
		if sameDay(dm.Date, date) {
			return dm.Meals, nil
		}
	}
	return nil, errors.New("no meals found for this date")
}

func (u *User) AddMealToDate(date time.Time, meal Meal) error {
	for i, dm := range u.DailyMenus {
		if sameDay(dm.Date, date) {
			u.DailyMenus[i].Meals = append(u.DailyMenus[i].Meals, meal)
			return nil
		}
	}

	u.DailyMenus = append(u.DailyMenus, DailyMenu{
		Date:  date,
		Meals: []Meal{meal},
	})
	return nil
}

func (u *User) RemoveMeal(date time.Time, mealIndex int) error {
	for i := range u.DailyMenus {
		if sameDay(u.DailyMenus[i].Date, date) {
			if mealIndex < 0 || mealIndex >= len(u.DailyMenus[i].Meals) {
				return errors.New("invalid meal index")
			}
			u.DailyMenus[i].Meals = append(u.DailyMenus[i].Meals[:mealIndex], u.DailyMenus[i].Meals[mealIndex+1:]...)
			return nil
		}
	}
	return errors.New("no meals found on this date")
}

func (u *User) GetDailyMacros(date time.Time) (float64, float64, float64, float64, error) {
	for _, dm := range u.DailyMenus {
		if sameDay(dm.Date, date) {
			return dm.GetDailyMacroSummary()
		}
	}
	return 0, 0, 0, 0, errors.New("no meals found on this date")
}

func (u *User) UpdateNutritionGoals() {
	if len(u.Measurements) == 0 {
		return
	}

	latest := u.Measurements[len(u.Measurements)-1]

	// Calcul du métabolisme de base (BMR) via la formule de Mifflin-St Jeor
	var bmr float64
	if u.Gender == Male {
		bmr = 10*latest.Weight + 6.25*latest.Height - 5*float64(u.Age) + 5
	} else {
		bmr = 10*latest.Weight + 6.25*latest.Height - 5*float64(u.Age) - 161
	}

	// Facteur d’activité par défaut (sédentaire)
	activityFactor := 1.5
	tdee := bmr * activityFactor

	// Ajustement selon l’objectif
	switch u.Goal {
	case WeightLoss:
		tdee -= 300
	case MuscleGain:
		tdee += 300
	case Maintenance:
		// pas de changement
	}

	u.CalorieNeeds = math.Round(tdee)

	// Calcul des macronutriments à partir du total calorique :
	// Protéines : 1.8 g/kg de poids
	// Lipides : 1 g/kg de poids
	// Glucides = reste des calories

	proteins := 1.8 * latest.Weight                       // g
	fats := 1.0 * latest.Weight                           // g
	proteinCals := proteins * 4                           // kcal
	fatCals := fats * 9                                   // kcal
	carbs := (u.CalorieNeeds - proteinCals - fatCals) / 4 // g

	u.ProteinNeeds = math.Round(proteins*100) / 100
	u.LipidNeeds = math.Round(fats*100) / 100
	u.CarohydratesNeeds = math.Round(carbs*100) / 100
}

func (u *User) GenerateNutritionChart(filename string) error {
	// Trie les menus par date
	sort.Slice(u.DailyMenus, func(i, j int) bool {
		return u.DailyMenus[i].Date.Before(u.DailyMenus[j].Date)
	})

	var dates []time.Time
	var calories, proteins, carbohydrates, lipids []float64

	for _, dm := range u.DailyMenus {
		dates = append(dates, dm.Date)
		cal, prot, carb, fat, _ := dm.GetDailyMacroSummary()
		calories = append(calories, cal)
		proteins = append(proteins, prot)
		carbohydrates = append(carbohydrates, carb)
		lipids = append(lipids, fat)
	}

	graph := chart.Chart{
		Width:  1024,
		Height: 512,
		XAxis: chart.XAxis{
			Name: "Date",
			ValueFormatter: func(v interface{}) string {
				if val, ok := v.(time.Time); ok {
					return val.Format("02/01")
				}
				return ""
			},
		},
		YAxis: chart.YAxis{
			Name: "Valeur (kcal / g)",
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Calories",
				XValues: dates,
				YValues: calories,
				Style:   chart.Style{StrokeColor: chart.GetDefaultColor(0), StrokeWidth: 2},
			},
			chart.TimeSeries{
				Name:    "Protéines (g)",
				XValues: dates,
				YValues: proteins,
				Style:   chart.Style{StrokeColor: chart.GetDefaultColor(1), StrokeWidth: 2},
			},
			chart.TimeSeries{
				Name:    "Glucides (g)",
				XValues: dates,
				YValues: carbohydrates,
				Style:   chart.Style{StrokeColor: chart.GetDefaultColor(2), StrokeWidth: 2},
			},
			chart.TimeSeries{
				Name:    "Lipides (g)",
				XValues: dates,
				YValues: lipids,
				Style:   chart.Style{StrokeColor: chart.GetDefaultColor(3), StrokeWidth: 2},
			},
		},
	}

	// Active la légende
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	// Création du fichier image
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return graph.Render(chart.PNG, f)
}

func (u *User) UpdateProfile(weight, height float64, age int, goal Goal, gender Gender, waist, neck, hip float64) error {
	if weight <= 0 || height <= 0 || age <= 0 {
		return errors.New("weight, height and age must be greater than 0")
	}

	u.Age = age
	u.Goal = goal
	u.Gender = gender

	bmi, err := CalculateBMI(weight, height)
	if err != nil {
		return err
	}

	bodyFat, err := CalculateBodyFat(gender, height, waist, neck, hip)
	if err != nil {
		return err
	}

	measurement := Measurement{
		Date:    time.Now(),
		Weight:  weight,
		Height:  height,
		BMI:     bmi,
		BodyFat: bodyFat,
	}

	u.Measurements = append(u.Measurements, measurement)
	return nil
}

func (u *User) GenerateBodyTrackingChart(filename string) error {
	if len(u.Measurements) == 0 {
		return errors.New("no measurements available to generate the chart")
	}

	// Trie les mesures par date
	sort.Slice(u.Measurements, func(i, j int) bool {
		return u.Measurements[i].Date.Before(u.Measurements[j].Date)
	})

	var dates []time.Time
	var bmiValues []float64
	var bodyFatValues []float64

	for _, m := range u.Measurements {
		dates = append(dates, m.Date)
		bmiValues = append(bmiValues, m.BMI)
		bodyFatValues = append(bodyFatValues, m.BodyFat)
	}

	graph := chart.Chart{
		Width:  1024,
		Height: 512,
		XAxis: chart.XAxis{
			Name: "Date",
			ValueFormatter: func(v interface{}) string {
				if val, ok := v.(time.Time); ok {
					return val.Format("02/01")
				}
				return ""
			},
		},
		YAxis: chart.YAxis{
			Name: "Value",
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "BMI",
				XValues: dates,
				YValues: bmiValues,
				Style:   chart.Style{StrokeColor: chart.GetDefaultColor(0), StrokeWidth: 2},
			},
			chart.TimeSeries{
				Name:    "Body Fat (%)",
				XValues: dates,
				YValues: bodyFatValues,
				Style:   chart.Style{StrokeColor: chart.GetDefaultColor(1), StrokeWidth: 2},
			},
		},
	}

	// Maintenant que graph est défini, on peut ajouter la légende
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return graph.Render(chart.PNG, f)
}
