package main

import (
	"fmt"
	"log"

	"github.com/zlahrouni/loganizer/internal/config"
)

func main_test() {
	fmt.Println("=== Test du package config ===")

	configs, err := config.LoadConfig("config.json")
	if err != nil {
		log.Printf("Erreur: %v", err)
		return
	}

	fmt.Printf("%d configurations chargées\n\n", len(configs))

	for i, cfg := range configs {
		fmt.Printf("Config %d:\n", i+1)
		fmt.Printf("  ID: %s\n", cfg.ID)
		fmt.Printf("  Path: %s\n", cfg.Path)
		fmt.Printf("  Type: %s\n", cfg.Type)
		fmt.Println()
	}
}
