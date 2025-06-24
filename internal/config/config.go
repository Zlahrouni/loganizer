package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type LogConfig struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

// Lit et valide le fichier de configuration JSON
func LoadConfig(configPath string) ([]LogConfig, error) {
	// Lire le fichier
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier %s: %v", configPath, err)
	}

	// Parser le JSON
	var configs []LogConfig
	err = json.Unmarshal(data, &configs)
	if err != nil {
		return nil, fmt.Errorf("format JSON invalide: %v", err)
	}

	// Valider les configurations
	err = validateConfigs(configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

// Vérifie que chaque configuration a tous les champs obligatoires
func validateConfigs(configs []LogConfig) error {
	for i, config := range configs {
		// Vérifier que l'ID n'est pas vide
		if strings.TrimSpace(config.ID) == "" {
			return fmt.Errorf("configuration %d: le champ 'id' est obligatoire", i+1)
		}

		// Vérifier que le path n'est pas vide
		if strings.TrimSpace(config.Path) == "" {
			return fmt.Errorf("configuration %d: le champ 'path' est obligatoire", i+1)
		}

		// Vérifier que le type n'est pas vide
		if strings.TrimSpace(config.Type) == "" {
			return fmt.Errorf("configuration %d: le champ 'type' est obligatoire", i+1)
		}
	}

	return nil
}
