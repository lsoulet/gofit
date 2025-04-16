package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lsoulet/gofit/db"
	"github.com/lsoulet/gofit/fdc"
	"github.com/lsoulet/gofit/models"
)

var (
	selectedUser models.User
	menuCreated  bool

	awaitingMealType bool
	mealTypeCallback func(string)

	awaitingMealDescription bool
	mealTypeSelected        models.MealType
	mealDescriptionCallback func(string)

	awaitingMealChoice bool
	mealChoiceCallback func(string)

	awaitingQuantity bool
	quantityCallback func(string)

	awaitingUserChoice bool
	userChoiceCallback func(string)

	awaitingDate bool
	dateCallback func(string)

	awaitingMenuChoice bool
	menuChoiceCallback func(string)

	awaitingFirstName bool
	firstNameCallback func(string)

	awaitingLastName bool
	lastNameCallback func(string)

	awaitingAge bool
	ageCallback func(string)

	awaitingGender bool
	genderCallback func(string)

	awaitingGoal bool
	goalCallback func(string)

	reader *bufio.Reader
)

// Command représente une commande saisie par l'utilisateur
type Command struct {
	Action string
	Args   []string
}

func startAsyncSaver(saveCh <-chan any) {
	for entity := range saveCh {
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("Erreur lors de la sauvegarde : timeout")
		default:
			if err := db.DB.Create(entity).Error; err != nil {
				fmt.Println("Erreur lors de la sauvegarde :", err)
			}
		}
	}
}

func startUserInputListener(commandChan chan<- Command) {
	reader = bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erreur lors de la lecture :", err)
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Étape 1 : on attend un type de repas
		if awaitingMealType && mealTypeCallback != nil {
			mealTypeCallback(input)
			awaitingMealType = false
			continue
		}

		// Étape 2 : on attend la description du repas
		if awaitingMealDescription && mealDescriptionCallback != nil {
			mealDescriptionCallback(input)
			awaitingMealDescription = false
			continue
		}

		// Étape 3 : on attend le choix d'un repas
		if awaitingMealChoice && mealChoiceCallback != nil {
			mealChoiceCallback(input)
			awaitingMealChoice = false
			continue
		}

		// Étape 4 : on attend une quantité
		if awaitingQuantity && quantityCallback != nil {
			quantityCallback(input)
			awaitingQuantity = false
			continue
		}

		// Étape 5 : on attend le choix d'un utilisateur
		if awaitingUserChoice && userChoiceCallback != nil {
			userChoiceCallback(input)
			awaitingUserChoice = false
			continue
		}

		// Étape 6 : on attend une date
		if awaitingDate && dateCallback != nil {
			dateCallback(input)
			awaitingDate = false
			continue
		}

		// Étape 7 : on attend le choix d'un menu
		if awaitingMenuChoice && menuChoiceCallback != nil {
			menuChoiceCallback(input)
			awaitingMenuChoice = false
			continue
		}

		// Étape 8 : on attend le prénom
		if awaitingFirstName && firstNameCallback != nil {
			firstNameCallback(input)
			awaitingFirstName = false
			continue
		}

		// Étape 9 : on attend le nom
		if awaitingLastName && lastNameCallback != nil {
			lastNameCallback(input)
			awaitingLastName = false
			continue
		}

		// Étape 10 : on attend l'âge
		if awaitingAge && ageCallback != nil {
			ageCallback(input)
			awaitingAge = false
			continue
		}

		// Étape 11 : on attend le genre
		if awaitingGender && genderCallback != nil {
			genderCallback(input)
			awaitingGender = false
			continue
		}

		// Étape 12 : on attend l'objectif
		if awaitingGoal && goalCallback != nil {
			goalCallback(input)
			awaitingGoal = false
			continue
		}

		// Commande classique
		parts := strings.Fields(input)

		// On ne vérifie le préfixe gofit que si on n'attend pas d'entrée spécifique
		if !awaitingMealType && !awaitingMealDescription && !awaitingQuantity && !awaitingMealChoice &&
			!awaitingUserChoice && !awaitingDate && !awaitingMenuChoice && !awaitingFirstName &&
			!awaitingLastName && !awaitingAge && !awaitingGender && !awaitingGoal {
			if !strings.HasPrefix(input, "gofit") {
				fmt.Println("Toutes les commandes doivent commencer par 'gofit'")
				continue
			}

			if len(parts) < 2 {
				fmt.Println("Commande incomplète.")
				continue
			}
		}

		command := Command{
			Action: parts[1],
			Args:   parts[2:],
		}

		commandChan <- command
	}
}

func handleCommand(cmd Command) bool {
	switch cmd.Action {
	case "search":
		if len(cmd.Args) < 1 {
			fmt.Println("Usage : gofit search <nom de l'aliment>")
			return false
		}
		results, err := fdc.SearchFood(cmd.Args[0])
		if err != nil {
			fmt.Println("Erreur lors de la recherche :", err)
			break
		}
		if len(results) == 0 {
			fmt.Println("Aucun résultat trouvé.")
			break
		}
		fmt.Println("Résultats trouvés :")
		for _, r := range results {
			fmt.Println("-", r)
		}

	case "detail":
		if len(cmd.Args) < 1 {
			fmt.Println("Usage : gofit detail <fdcId>")
			return false
		}
		id, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			fmt.Println("fdcId invalide :", cmd.Args[0])
			return false
		}
		name, calories, proteins, carbs, lipids, err := fdc.GetFoodDetails(id)
		if err != nil {
			fmt.Println("Erreur lors de la récupération :", err)
			break
		}
		fmt.Println("Détails nutritionnels :")
		fmt.Printf("Nom : %s\n", name)
		fmt.Printf("Calories : %.2f kcal\n", calories)
		fmt.Printf("Protéines : %.2f g\n", proteins)
		fmt.Printf("Glucides : %.2f g\n", carbs)
		fmt.Printf("Lipides : %.2f g\n", lipids)
		fmt.Printf("Quantité : %.2f g\n", 100.0)

	case "addfood":
		if len(cmd.Args) < 1 {
			fmt.Println("Usage : gofit addfood <fdcId>")
			return false
		}

		fdcID, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			fmt.Println("fdcId invalide :", cmd.Args[0])
			return false
		}

		// Récupérer les détails de l'aliment
		name, _, _, _, _, err := fdc.GetFoodDetails(fdcID)
		if err != nil {
			fmt.Println("Erreur lors de la récupération de l'aliment :", err)
			return false
		}

		// Récupérer la liste des repas
		meals, err := fdc.GetMeals()
		if err != nil {
			fmt.Println("Erreur lors de la récupération des repas :", err)
			return false
		}

		if len(meals) == 0 {
			fmt.Println("Aucun repas n'a été créé. Veuillez d'abord créer un repas avec 'gofit addmeal'.")
			return false
		}

		// Afficher la liste des repas
		fmt.Printf("\nAliment sélectionné : %s\n\n", name)
		fmt.Println("Choisissez le repas auquel ajouter cet aliment :")
		for i, meal := range meals {
			fmt.Printf("%d. %s (%s)\n", i+1, meal.Description, meal.Type)
		}

		var selectedMeal models.Meal
		awaitingMealChoice = true
		mealChoiceCallback = func(choiceStr string) {
			// Convertir le choix en nombre
			choice, err := strconv.Atoi(strings.TrimSpace(choiceStr))
			if err != nil || choice < 1 || choice > len(meals) {
				fmt.Printf("Choix invalide. Veuillez entrer un nombre entre 1 et %d.\n", len(meals))
				awaitingMealChoice = true
				return
			}

			// Stocker le repas sélectionné
			selectedMeal = meals[choice-1]

			// Demander la quantité
			fmt.Println("\nVeuillez saisir la quantité en grammes :")
			awaitingQuantity = true
			quantityCallback = func(quantityStr string) {
				quantity, err := strconv.ParseFloat(strings.TrimSpace(quantityStr), 64)
				if err != nil || quantity <= 0 {
					fmt.Println("Quantité invalide. Veuillez entrer un nombre positif.")
					awaitingQuantity = true
					return
				}

				// Ajouter l'aliment au repas choisi
				if err := fdc.AddFoodToMeal(selectedMeal.ID, fdcID, quantity); err != nil {
					fmt.Println("Erreur lors de l'ajout de l'aliment au repas :", err)
					return
				}
				fmt.Printf("\n✅ %.0fg de %s ajoutés au repas '%s'\n", quantity, name, selectedMeal.Description)
			}
		}

	case "addmeal":
		// Récupérer la liste des repas existants
		meals, err := fdc.GetMeals()
		if err != nil {
			fmt.Println("Erreur lors de la récupération des repas :", err)
			return false
		}

		// Récupérer la liste des menus
		menus, err := fdc.GetDailyMenus()
		if err != nil {
			fmt.Println("Erreur lors de la récupération des menus :", err)
			return false
		}

		if len(menus) == 0 {
			fmt.Println("Aucun menu journalier enregistré. Veuillez d'abord créer un menu avec 'gofit addday'.")
			return false
		}

		// Afficher la liste des menus
		fmt.Println("\nChoisissez le menu auquel ajouter ce repas :")
		for i, menu := range menus {
			fmt.Printf("%d. %s %s - %s\n", i+1, menu.User.FirstName, menu.User.LastName, menu.Date.Format("02/01/2006"))
		}

		var selectedMenu models.DailyMenu
		awaitingMenuChoice = true
		menuChoiceCallback = func(input string) {
			choice, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil || choice < 1 || choice > len(menus) {
				fmt.Println("Choix invalide. Veuillez entrer un nombre entre 1 et", len(menus))
				awaitingMenuChoice = true
				return
			}

			selectedMenu = menus[choice-1]

			// Afficher les repas existants
			fmt.Println("\nChoisissez un repas existant :")
			for i, meal := range meals {
				fmt.Printf("%d. %s (%s)\n", i+1, meal.Description, meal.Type)
			}

			awaitingMealChoice = true
			mealChoiceCallback = func(input string) {
				choice, err := strconv.Atoi(strings.TrimSpace(input))
				if err != nil || choice < 1 || choice > len(meals) {
					fmt.Printf("Choix invalide. Veuillez entrer un nombre entre 1 et %d.\n", len(meals))
					awaitingMealChoice = true
					return
				}

				// Utiliser un repas existant
				selectedMeal := meals[choice-1]

				// Créer un nouveau repas avec les mêmes caractéristiques
				if err := fdc.AddMealToDailyMenu(selectedMenu.ID, selectedMeal.Type, selectedMeal.Description); err != nil {
					fmt.Println("Erreur lors de la création du repas :", err)
					return
				}

				fmt.Printf("\n✅ Repas '%s' (%s) ajouté au menu de %s %s le %s\n",
					selectedMeal.Description, selectedMeal.Type, selectedMenu.User.FirstName, selectedMenu.User.LastName, selectedMenu.Date.Format("02/01/2006"))
			}
		}

	case "newmeal":
		fmt.Println("Quel type de repas souhaitez-vous ajouter ?")
		mealOptions := []models.MealType{
			models.Breakfast,
			models.Lunch,
			models.Dinner,
			models.Snack,
		}

		for i, mt := range mealOptions {
			fmt.Printf("%d. %s\n", i+1, mt)
		}

		fmt.Print("Tape le numéro correspondant : ")

		awaitingMealType = true
		mealTypeCallback = func(input string) {
			choice, err := strconv.Atoi(input)
			if err != nil || choice < 1 || choice > len(mealOptions) {
				fmt.Println("Entrée invalide. Tape un numéro entre 1 et", len(mealOptions))
				awaitingMealType = true
				return
			}

			mealTypeSelected = mealOptions[choice-1]

			fmt.Println("Saisis maintenant une description pour ce repas (ex: \"Déjeuner du mardi\") :")
			awaitingMealDescription = true
			mealDescriptionCallback = func(desc string) {
				meal := models.Meal{
					Type:        mealTypeSelected,
					Description: desc,
				}

				if err := fdc.AddMeal(meal); err != nil {
					fmt.Println("Erreur lors de la sauvegarde du repas :", err)
				} else {
					fmt.Printf("✔ Repas '%s' (%s) ajouté avec succès !\n", desc, mealTypeSelected)
				}
			}
		}

	case "report":
		fdc.ListMeals()
		fmt.Println("Génération du bilan nutritionnel journalier...")
		if err := fdc.GenerateNutritionalReport(); err != nil {
			fmt.Println("Erreur lors de la génération du rapport :", err)
		}

	case "adduser":
		// Demander le prénom
		fmt.Println("\nEntrez le prénom de l'utilisateur :")
		awaitingFirstName = true
		var firstName string
		firstNameCallback = func(input string) {
			firstName = strings.TrimSpace(input)
			if firstName == "" {
				fmt.Println("Le prénom ne peut pas être vide.")
				awaitingFirstName = true
				return
			}

			// Demander le nom
			fmt.Println("\nEntrez le nom de l'utilisateur :")
			awaitingLastName = true
			lastNameCallback = func(input string) {
				lastName := strings.TrimSpace(input)
				if lastName == "" {
					fmt.Println("Le nom ne peut pas être vide.")
					awaitingLastName = true
					return
				}

				// Demander l'âge
				fmt.Println("\nEntrez l'âge de l'utilisateur :")
				awaitingAge = true
				ageCallback = func(input string) {
					age, err := strconv.Atoi(strings.TrimSpace(input))
					if err != nil || age <= 0 {
						fmt.Println("L'âge doit être un nombre positif.")
						awaitingAge = true
						return
					}

					// Demander le genre
					fmt.Println("\nChoisissez le genre :")
					fmt.Println("1. Homme")
					fmt.Println("2. Femme")
					awaitingGender = true
					genderCallback = func(input string) {
						choice, err := strconv.Atoi(strings.TrimSpace(input))
						if err != nil || choice < 1 || choice > 2 {
							fmt.Println("Choix invalide. Veuillez entrer 1 ou 2.")
							awaitingGender = true
							return
						}

						var gender models.Gender
						if choice == 1 {
							gender = models.Male
						} else {
							gender = models.Female
						}

						// Demander l'objectif
						fmt.Println("\nChoisissez l'objectif :")
						fmt.Println("1. Perte de poids")
						fmt.Println("2. Maintien")
						fmt.Println("3. Prise de masse")
						awaitingGoal = true
						goalCallback = func(input string) {
							choice, err := strconv.Atoi(strings.TrimSpace(input))
							if err != nil || choice < 1 || choice > 3 {
								fmt.Println("Choix invalide. Veuillez entrer un nombre entre 1 et 3.")
								awaitingGoal = true
								return
							}

							var goal models.Goal
							switch choice {
							case 1:
								goal = models.WeightLoss
							case 2:
								goal = models.Maintenance
							case 3:
								goal = models.MuscleGain
							}

							// Créer l'utilisateur
							if err := fdc.CreateUser(firstName, lastName, age, gender, goal); err != nil {
								fmt.Println("Erreur lors de la création de l'utilisateur :", err)
								return
							}

							fmt.Printf("\n✅ Utilisateur %s %s créé avec succès !\n", firstName, lastName)
						}
					}
				}
			}
		}

	case "addmenu":
		// Récupérer la liste des utilisateurs
		users, err := fdc.GetUsers()
		if err != nil {
			fmt.Println("Erreur lors de la récupération des utilisateurs :", err)
			return false
		}

		if len(users) == 0 {
			fmt.Println("Aucun utilisateur enregistré. Veuillez d'abord créer un utilisateur.")
			return false
		}

		// Réinitialiser l'état
		menuCreated = false

		// Afficher la liste des utilisateurs
		fmt.Println("\nChoisissez l'utilisateur pour ce menu :")
		for i, user := range users {
			fmt.Printf("%d. %s %s\n", i+1, user.FirstName, user.LastName)
		}

		awaitingUserChoice = true
		userChoiceCallback = func(choiceStr string) {
			// Convertir le choix en nombre
			choice, err := strconv.Atoi(strings.TrimSpace(choiceStr))
			if err != nil || choice < 1 || choice > len(users) {
				fmt.Printf("Choix invalide. Veuillez entrer un nombre entre 1 et %d.\n", len(users))
				awaitingUserChoice = true
				return
			}

			// Stocker l'utilisateur sélectionné
			selectedUser = users[choice-1]

			// Demander la date
			fmt.Println("\nEntrez la date au format JJ/MM/AAAA :")
			awaitingDate = true
			dateCallback = func(dateStr string) {
				// Parser la date
				date, err := time.Parse("02/01/2006", strings.TrimSpace(dateStr))
				if err != nil {
					fmt.Println("Format de date invalide. Utilisez le format JJ/MM/AAAA.")
					awaitingDate = true
					return
				}

				// Créer le menu journalier
				if err := fdc.CreateDailyMenu(selectedUser.ID, date); err != nil {
					fmt.Println("Erreur lors de la création du menu :", err)
					return
				}

				fmt.Printf("\n✅ Menu journalier créé pour %s %s le %s\n",
					selectedUser.FirstName, selectedUser.LastName, date.Format("02/01/2006"))
				menuCreated = true
			}
		}

		// Attendre que le menu soit créé
		if !menuCreated {
			return false
		}
		return true

	default:
		fmt.Println("Commande inconnue :", cmd.Action)
	}
	return false
}

var saveChan chan any

func main() {
	db.InitDatabase()

	saveChan = make(chan any, 100)
	go startAsyncSaver(saveChan)

	commandChan := make(chan Command)
	go startUserInputListener(commandChan)

	for {
		select {
		case cmd := <-commandChan:
			if stop := handleCommand(cmd); stop {
				return
			}
		}
	}
}
