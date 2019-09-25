/*
 * This file is part of GCSC.
 *
 * GCSC is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * GCSC is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with GCSC. If not, see <https://www.gnu.org/licenses/>.
 */

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
