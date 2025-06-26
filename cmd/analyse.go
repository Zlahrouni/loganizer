package cmd

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zlahrouni/loganizer/internal/analyzer"
	"github.com/zlahrouni/loganizer/internal/config"
)

var (
	configFile   string
	outputFile   string
	showStatus   bool
	sortBy       string
	filterType   string
	filterStatus string
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyse les fichiers de logs selon la configuration",
	Long: `La commande analyze traite les fichiers de logs spécifiés dans le fichier
de configuration en utilisant la concurrence pour optimiser les performances.

Elle génère des statistiques détaillées sur chaque fichier analysé et peut
exporter les résultats dans différents formats.`,
	Example: `  # Analyse basique avec fichier de configuration
  loganizer analyze -c config.json

  # Analyse avec sauvegarde des résultats
  loganizer analyze -c config.json -o results.txt

  # Analyse avec statut temps réel et tri par nombre d'entrées
  loganizer analyze -c config.json --status --sort entries

  # Analyse avec filtrage par type
  loganizer analyze -c config.json --filter-type nginx-access`,
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&configFile, "config", "c", "", "fichier de configuration JSON (obligatoire)")
	analyzeCmd.MarkFlagRequired("config")

	analyzeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "fichier de sortie pour les résultats")
	analyzeCmd.Flags().BoolVar(&showStatus, "status", false, "affiche le statut en temps réel pendant l'analyse")

	analyzeCmd.Flags().StringVar(&sortBy, "sort", "", "trier par: 'id', 'entries', 'duration', 'status' (optionnel)")
	analyzeCmd.Flags().StringVar(&filterType, "filter-type", "", "filtrer par type de log (optionnel)")
	analyzeCmd.Flags().StringVar(&filterStatus, "filter-status", "", "filtrer par statut: 'success', 'error' (optionnel)")
}

type AnalysisStatus struct {
	ID       string
	Path     string
	Type     string
	Status   string
	Entries  int
	Duration time.Duration
	Error    string
}

func runAnalyze(cmd *cobra.Command, args []string) error {

	fmt.Printf("Démarrage de l'analyse avec la configuration: %s\n\n", configFile)

	if showStatus {
		fmt.Println("Chargement de la configuration...")
	}

	configs, err := config.LoadConfig(configFile)
	if err != nil {
		return fmt.Errorf("erreur lors du chargement de la configuration: %v", err)
	}

	if showStatus {
		fmt.Printf("%d configurations chargées\n\n", len(configs))
	}

	var filePaths []string
	configMap := make(map[string]config.LogConfig)

	for _, cfg := range configs {
		filePaths = append(filePaths, cfg.Path)
		configMap[cfg.Path] = cfg
	}

	if showStatus {
		fmt.Println("Analyse en cours...")
		fmt.Println("────────────────────────────────────")
	}

	start := time.Now()
	results, errs := analyzer.AnalyzeLogs(filePaths)
	totalDuration := time.Since(start)

	statuses := createAnalysisStatuses(configs, results, errs, configMap)

	statuses = applyFilters(statuses)

	sortResults(statuses)

	displayResults(statuses, totalDuration)

	if outputFile != "" {
		err := saveResults(statuses, totalDuration)
		if err != nil {
			return fmt.Errorf("erreur lors de la sauvegarde: %v", err)
		}
		fmt.Printf("\nRésultats sauvegardés dans: %s\n", outputFile)
	}

	return nil
}

func createAnalysisStatuses(configs []config.LogConfig, results []analyzer.AnalysisResult, errs []error, configMap map[string]config.LogConfig) []AnalysisStatus {
	var statuses []AnalysisStatus

	resultMap := make(map[string]analyzer.AnalysisResult)
	for _, result := range results {
		resultMap[result.FilePath] = result
	}

	errorMap := make(map[string]error)
	for _, err := range errs {
		var notFoundErr *analyzer.FileNotFoundOrUnreadableError
		var parsingErr *analyzer.ParsingError

		switch {
		case errors.As(err, &notFoundErr):
			errorMap[notFoundErr.FilePath] = err
		case errors.As(err, &parsingErr):
			errorMap[parsingErr.FilePath] = err
		}
	}

	for _, cfg := range configs {
		status := AnalysisStatus{
			ID:   cfg.ID,
			Path: cfg.Path,
			Type: cfg.Type,
		}

		if result, exists := resultMap[cfg.Path]; exists {
			status.Status = "Succès"
			status.Entries = result.Entries
			status.Duration = result.Duration
		} else if err, exists := errorMap[cfg.Path]; exists {
			status.Status = "Erreur"
			status.Error = err.Error()
		} else {
			status.Status = "Non traité"
		}

		statuses = append(statuses, status)
	}

	return statuses
}

func applyFilters(statuses []AnalysisStatus) []AnalysisStatus {
	var filtered []AnalysisStatus

	for _, status := range statuses {

		if filterType != "" && !strings.Contains(strings.ToLower(status.Type), strings.ToLower(filterType)) {
			continue
		}

		if filterStatus != "" {
			switch strings.ToLower(filterStatus) {
			case "success":
				if !strings.Contains(status.Status, "Succès") {
					continue
				}
			case "error":
				if !strings.Contains(status.Status, "Erreur") {
					continue
				}
			}
		}

		filtered = append(filtered, status)
	}

	return filtered
}

func sortResults(statuses []AnalysisStatus) {
	switch strings.ToLower(sortBy) {
	case "id":
		sort.Slice(statuses, func(i, j int) bool {
			return statuses[i].ID < statuses[j].ID
		})
	case "entries":
		sort.Slice(statuses, func(i, j int) bool {
			return statuses[i].Entries > statuses[j].Entries
		})
	case "duration":
		sort.Slice(statuses, func(i, j int) bool {
			return statuses[i].Duration > statuses[j].Duration
		})
	case "status":
		sort.Slice(statuses, func(i, j int) bool {
			return statuses[i].Status < statuses[j].Status
		})
	}
}

func displayResults(statuses []AnalysisStatus, totalDuration time.Duration) {
	fmt.Println("\nRÉSULTATS D'ANALYSE")
	fmt.Println("════════════════════════════════════════════════════════════════")

	successCount := 0
	totalEntries := 0

	for _, status := range statuses {
		fmt.Printf("ID: %s\n", status.ID)
		fmt.Printf("   Chemin: %s\n", status.Path)
		fmt.Printf("   Type: %s\n", status.Type)
		fmt.Printf("   Statut: %s\n", status.Status)

		if status.Entries > 0 {
			fmt.Printf("   Entrées: %d\n", status.Entries)
			fmt.Printf("   Durée: %v\n", status.Duration)
			successCount++
			totalEntries += status.Entries
		}

		if status.Error != "" {
			fmt.Printf("   Erreur: %s\n", status.Error)
		}

		fmt.Println()
	}

	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Printf("STATISTIQUES GLOBALES\n")
	fmt.Printf("   Fichiers traités avec succès: %d/%d\n", successCount, len(statuses))
	fmt.Printf("   Total des entrées analysées: %d\n", totalEntries)
	fmt.Printf("   Durée totale d'analyse: %v\n", totalDuration)

	if successCount > 0 {
		avgEntries := float64(totalEntries) / float64(successCount)
		fmt.Printf("   Moyenne d'entrées par fichier: %.1f\n", avgEntries)
	}
}

func saveResults(statuses []AnalysisStatus, totalDuration time.Duration) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "LOGANIZER - RAPPORT D'ANALYSE\n")
	fmt.Fprintf(file, "Généré le: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "Durée totale: %v\n\n", totalDuration)

	for _, status := range statuses {
		fmt.Fprintf(file, "ID: %s\n", status.ID)
		fmt.Fprintf(file, "Chemin: %s\n", status.Path)
		fmt.Fprintf(file, "Type: %s\n", status.Type)
		fmt.Fprintf(file, "Statut: %s\n", status.Status)

		if status.Entries > 0 {
			fmt.Fprintf(file, "Entrées: %d\n", status.Entries)
			fmt.Fprintf(file, "Durée: %v\n", status.Duration)
		}

		if status.Error != "" {
			fmt.Fprintf(file, "Erreur: %s\n", status.Error)
		}

		fmt.Fprintf(file, "\n")
	}

	return nil
}
