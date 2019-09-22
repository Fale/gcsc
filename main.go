package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fale/gcsc/pkg/clean"
)

func main() {
	rootCmd := initCmd()

	if err := config(rootCmd); err != nil {
		fmt.Printf("an error occurred during configuration parsing: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cleanFn(cmd *cobra.Command, args []string) error {
	p := clean.Parameters{
		ProjectID:          viper.GetString("project-id"),
		CleanAutoBackups:   viper.GetBool("automatic"),
		CleanManualBackups: viper.GetBool("manual"),
		DryRun:             viper.GetBool("dry-run"),
	}
	err := viper.UnmarshalKey("retention-policies", &p.RetentionPolicies)
	if err != nil {
		panic(err)
	}
	return clean.Execute(p)
}
