package main

import "github.com/spf13/cobra"

func initCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use: "snapshot-cleaner",
	}

	rootCmd.PersistentFlags().StringP("project-id", "p", "", "Google Cloud Project ID")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Dry run mode")
	rootCmd.PersistentFlags().Bool("automatic", true, "Include automatic backups")
	rootCmd.PersistentFlags().Bool("manual", false, "Include manual backups")

	var pubsubCmd = &cobra.Command{
		Use:   "pubsub",
		Short: "listen to a pub/sub stream",
		//    RunE:  serve,
	}
	rootCmd.AddCommand(pubsubCmd)

	var cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "execute a cleaning",
		RunE:  clean,
	}

	rootCmd.AddCommand(cleanCmd)

	return rootCmd
}
