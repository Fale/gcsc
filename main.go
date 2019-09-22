package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "gcsc",
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

	var httpCmd = &cobra.Command{
		Use:   "http",
		Short: "listen to HTTP port",
		RunE:  httpFn,
	}
	rootCmd.AddCommand(httpCmd)

	if err := config(rootCmd); err != nil {
		fmt.Printf("an error occurred during configuration parsing: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
