package cmd

import (
    "github.com/spf13/cobra"
    "github.com/zlahrouni/loganizer/internal/reporter"
)

var addLogCmd = &cobra.Command{
    Use:   "add-log",
    Short: "Ajoute un log à la config",
    RunE: func(cmd *cobra.Command, args []string) error {
        id, _ := cmd.Flags().GetString("id")
        path, _ := cmd.Flags().GetString("path")
        logType, _ := cmd.Flags().GetString("type")
        file, _ := cmd.Flags().GetString("file")
        
        return reporter.AddLog(file, id, path, logType)
    },
}

func init() {
    rootCmd.AddCommand(addLogCmd)
    addLogCmd.Flags().String("id", "", "ID du log")
    addLogCmd.Flags().String("path", "", "Chemin du fichier")
    addLogCmd.Flags().String("type", "", "Type de log")
    addLogCmd.Flags().String("file", "", "Fichier de config")
    addLogCmd.MarkFlagRequired("id")
    addLogCmd.MarkFlagRequired("path")
    addLogCmd.MarkFlagRequired("type")
    addLogCmd.MarkFlagRequired("file")
}