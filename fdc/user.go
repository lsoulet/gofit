package fdc

import (
	"fmt"

	"github.com/lsoulet/gofit/db"
	"github.com/lsoulet/gofit/models"
)

func CreateUser(firstName, lastName string, age int, gender models.Gender, goal models.Goal) error {
	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
		Gender:    gender,
		Goal:      goal,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return fmt.Errorf("erreur lors de la création de l'utilisateur : %w", err)
	}

	return nil
}

func ListUsers() error {
	users, err := GetUsers()
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("Aucun utilisateur enregistré.")
		return nil
	}

	fmt.Println("👤 Utilisateurs enregistrés :")
	for i, user := range users {
		fmt.Printf("%d. %s %s (%d ans)\n", i+1, user.FirstName, user.LastName, user.Age)
	}
	return nil
}
