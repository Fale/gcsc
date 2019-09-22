package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fale/gcsc/pkg/clean"
)

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
