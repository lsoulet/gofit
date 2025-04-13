package models

type Meal struct {
	Type        MealType
	FoodItems   []FoodItem
	Description string
}

func (m *Meal) CalculateMacros() (float64, float64, float64, float64) {
	var calories, proteins, carbohydrates, lipids float64
	for _, food := range m.FoodItems {
		calories += food.Calories
		proteins += food.Proteins
		carbohydrates += food.Carbohydrates
		lipids += food.Lipids
	}
	return calories, proteins, carbohydrates, lipids
}
