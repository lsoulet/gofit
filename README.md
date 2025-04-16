# GoFit - Assistant de Suivi Nutritionnel

GoFit est une application en ligne de commande (CLI) dédiée au suivi nutritionnel, permettant aux utilisateurs de suivre leur alimentation quotidienne et d'atteindre leurs objectifs de santé grâce à l'API FoodData Central (FDC).

## Fonctionnalités

- Recherche d'aliments dans la base de données FoodData Central
- Gestion des repas quotidiens avec portions personnalisées
- Génération de rapports nutritionnels détaillés
- Calcul automatique des macronutriments et calories
- Visualisation des données nutritionnelles sous forme de tableaux
- Historique des repas consommés
- Suggestions de menus quotidiens équilibrés

## Prérequis

- Go 1.23 ou supérieur
- PostgreSQL 14 ou supérieur
- Une clé API FoodData Central (obtenue sur [https://fdc.nal.usda.gov/api-key-signup.html](https://fdc.nal.usda.gov/api-key-signup.html))

## Installation

1. Cloner le dépôt :
```bash
git clone https://github.com/lsoulet/gofit.git
cd gofit
```

2. Installer les dépendances :
```bash
go mod download
```

3. Configurer la base de données PostgreSQL :
```bash
psql -U postgres
CREATE DATABASE gofit;
```

4. Configurer les variables d'environnement dans un fichier `.env` :
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=votre_mot_de_passe
DB_NAME=gofit
FDC_API_KEY=votre_clé_api_fdc
```

## Utilisation

1. Lancer l'application :
```bash
go run main.go
```

2. Commandes disponibles :

**Note** : Toutes les commandes doivent être préfixées par `gofit`

### Gestion des aliments
- `search [terme]` : Rechercher un aliment dans la base FDC
  ```bash
  gofit search pomme
  ```

- `detail [fdc_id]` : Voir les détails nutritionnels d'un aliment
  ```bash
  gofit detail 173939
  ```

- `addfood [fdc_id]` : Ajouter un aliment à un repas existant
  ```bash
  gofit addfood 173939
  ```

### Gestion des repas
- `newmeal` : Créer un nouveau repas type
  ```bash
  gofit newmeal
  ```

- `addmeal` : Ajouter un repas existant à un menu journalier
  ```bash
  gofit addmeal
  ```

### Gestion des utilisateurs
- `adduser` : Créer un nouvel utilisateur
  ```bash
  gofit adduser
  ```

### Menus journaliers
- `addmenu` : Créer un nouveau menu journalier
  ```bash
  gofit addmenu
  ```

### Rapports
- `report` : Générer un rapport nutritionnel pour tous les repas
  ```bash
  gofit report
  ```

## Structure du projet

```
gofit/
├── main.go           # Point d'entrée de l'application
├── models/          # Modèles de données
│   ├── meal.go
│   └── ...
├── fdc/             # Intégration avec l'API FoodData Central
│   ├── client.go
│   ├── daily_menu.go
│   └── report.go
├── db/              # Gestion de la base de données
│   └── ...
└── cmd/             # Commandes CLI
    └── ...
```

## Base de données

L'application utilise PostgreSQL pour stocker :
- Les repas enregistrés
- Les menus types
- L'historique des rapports nutritionnels
- Les préférences utilisateur

## Développement

### Contribution

1. Forker le projet
2. Créer une branche pour votre fonctionnalité
3. Commiter vos changements
4. Pousser vers la branche
5. Créer une Pull Request

## Support

Pour toute question ou problème :
1. Consulter les [Issues](https://github.com/lsoulet/gofit/issues)
2. Créer une nouvelle issue si nécessaire


- `main.go` : Point d'entrée de l'application
- `models/` : Définitions des structures de données
- `services/` : Logique métier (à venir)
- `storage/` : Gestion de la persistance des données (à venir)
