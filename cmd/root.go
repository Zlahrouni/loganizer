package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loganizer",
	Short: "Un analyseur de logs concurrent et efficace",
	Long: `Loganizer est un outil d'analyse de logs qui permet de traiter
plusieurs fichiers de logs en parallèle avec différents formats.

Il utilise la concurrence pour analyser rapidement de gros volumes
de logs et fournit des statistiques détaillées sur chaque fichier.`,
	Example: `  # Analyser avec un fichier de configuration
  loganizer analyze -c config.json

  # Analyser avec sortie dans un fichier
  loganizer analyze -c config.json -o results.txt

  # Analyser avec affichage du statut en temps réel
  loganizer analyze -c config.json --status`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de l'exécution de la commande: %v\n", err)
		os.Exit(1)
	}
}

func init() {
}
