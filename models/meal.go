package models

type Meal struct {
	ID            uint `gorm:"primaryKey"`
	DailyMenuID   uint
	Type          MealType
	Description   string
	Calories      float64
	Proteins      float64
	Carbohydrates float64
	Lipids        float64
}

func (m *Meal) GetMacros() (float64, float64, float64, float64) {
	return m.Calories, m.Proteins, m.Carbohydrates, m.Lipids
}
