package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "snapshot-cleaner",
	}

	rootCmd.PersistentFlags().StringP("project-id", "p", "", "Google Cloud Project ID")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Dry run mode")
	rootCmd.PersistentFlags().Bool("automatic", true, "Include automatic backups")
	rootCmd.PersistentFlags().Bool("manual", false, "Include manual backups")

	var cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "execute a cleaning",
		RunE:  cleanFn,
	}

	rootCmd.AddCommand(cleanCmd)

	if err := config(rootCmd); err != nil {
		fmt.Printf("an error occurred during configuration parsing: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
