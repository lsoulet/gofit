package models

import (
	"errors"
	"math"
	"time"
)

type Measurement struct {
	Date    time.Time
	Weight  float64
	Height  float64
	BodyFat float64
	BMI     float64
}

func CalculateBMI(weight, height float64) (float64, error) {
	if height <= 0 {
		return 0, errors.New("height must be greater than 0")
	}
	heightInMeters := height / 100
	bmi := weight / (heightInMeters * heightInMeters)
	return math.Round(bmi*100) / 100, nil
}

func CalculateBodyFat(gender Gender, height, waistCircumference, neckCircumference, hipCircumference float64) (float64, error) {
	if height <= 0 {
		return 0, errors.New("height must be greater than 0")
	}
	if waistCircumference <= 0 || neckCircumference <= 0 {
		return 0, errors.New("waist and neck measurements must be positive")
	}

	var bodyFat float64

	switch gender {
	case Male:
		diff := waistCircumference - neckCircumference
		if diff <= 0 {
			return 0, errors.New("waist circumference must be greater than neck circumference for males")
		}
		bodyFat = math.Round((495 / (1.0324 - 0.19077*math.Log10(diff) + 0.15456*math.Log10(height))) - 450)
	case Female:
		if hipCircumference <= 0 {
			return 0, errors.New("hip circumference must be positive for females")
		}
		sum := waistCircumference + hipCircumference - neckCircumference
		if sum <= 0 {
			return 0, errors.New("the sum of waist + hip - neck circumference must be positive for females")
		}
		bodyFat = math.Round((495 / (1.29579 - 0.35004*math.Log10(sum) + 0.22100*math.Log10(height))) - 450)
	default:
		return 0, errors.New("invalid gender: must be 'male' or 'female'")
	}

	return bodyFat, nil
}
