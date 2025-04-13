package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lsoulet/gofit/fdc"
)

// Command représente une commande saisie par l'utilisateur
type Command struct {
	Action string
	Args   []string
}

func startUserInputListener(commandChan chan<- Command) {
	reader := bufio.NewReader(os.Stdin)

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

		parts := strings.Fields(input)
		if len(parts) < 1 || parts[0] != "gofit" {
			fmt.Println("Toutes les commandes doivent commencer par 'gofit'")
			continue
		}

		if len(parts) < 2 {
			fmt.Println("Commande incomplète.")
			continue
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
		results, err := fdc.SearchFood(cmd.Args[0])
		if err != nil {
			fmt.Println("Erreur lors de la recherche :", err)
			break
		}
		fmt.Println("Résultats trouvés :")
		for _, r := range results {
			fmt.Println("-", r)
		}

	case "detail":
		id, _ := strconv.Atoi(cmd.Args[0])
		nutrients, err := fdc.GetFoodDetails(id)
		if err != nil {
			fmt.Println("Erreur lors de la récupération :", err)
			break
		}
		fmt.Println("Détails nutritionnels :")
		for name, val := range nutrients {
			fmt.Printf("%s : %.2f\n", name, val)
		}

	case "add":
		if len(cmd.Args) < 2 {
			fmt.Println("Usage : gofit add <aliment> <quantité>")
			return false
		}
		// Implémentation de l'ajout
		fmt.Printf("Ajout de %s (quantité %s)\n", cmd.Args[0], cmd.Args[1])

	case "report":
		// Affichage du récapitulatif
		fmt.Println("Génération du bilan nutritionnel journalier...")

	case "exit":
		fmt.Println("Fermeture de GoFit...")
		return true

	default:
		fmt.Println("Commande inconnue :", cmd.Action)
	}

	return false
}

func main() {
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
