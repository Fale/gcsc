package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/fale/gcsc/pkg/clean"
)

func httpFn(cmd *cobra.Command, args []string) error {
	http.HandleFunc("/", handler)
	return http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := clean.Parameters{
		ProjectID:          viper.GetString("project-id"),
		CleanAutoBackups:   viper.GetBool("automatic"),
		CleanManualBackups: viper.GetBool("manual"),
		DryRun:             viper.GetBool("dry-run"),
	}
	err := viper.UnmarshalKey("retention-policies", &p.RetentionPolicies)
	if err != nil {
		fmt.Fprintf(w, "an error occurred: %v", err)
	}
	err = clean.Execute(p)
	if err != nil {
		fmt.Fprintf(w, "an error occurred: %v", err)
	}
}
