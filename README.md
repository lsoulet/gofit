# GoFit - Assistant de Suivi Nutritionnel

GoFit est une application en ligne de commande (CLI) dédiée au suivi nutritionnel, permettant aux utilisateurs de suivre leur alimentation et leur condition physique.

## Fonctionnalités

- Gestion du profil utilisateur (informations personnelles, objectifs)
- Suivi de la composition corporelle (IMC, masse grasse)
- Recherche d'aliments dans une base de données nutritionnelle
- Gestion des repas et des journées types
- Suivi des macronutriments quotidiens
- Génération de bilans nutritionnels

## Installation

```bash
git clone https://github.com/lsoulet/gofit.git
cd gofit
go build
```

## Utilisation

Lancer l'application :
```bash
./gofit
```

### Commandes disponibles

- `search [aliment]` : Rechercher un aliment dans la base de données
- `add [aliment] [quantité]` : Ajouter un aliment consommé
- `report` : Afficher le bilan nutritionnel
- `exit` : Quitter l'application

## Structure du projet

- `main.go` : Point d'entrée de l'application
- `models/` : Définitions des structures de données
- `services/` : Logique métier (à venir)
- `storage/` : Gestion de la persistance des données (à venir)

## Développement

Ce projet est en cours de développement. Les prochaines étapes incluent :
- Implémentation de la base de données d'aliments
- Ajout de la persistance des données
- Amélioration des calculs nutritionnels
- Extension des fonctionnalités de reporting
