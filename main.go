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
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "gcsc",
	}

	rootCmd.PersistentFlags().StringP("project-id", "p", "", "Google Cloud Project ID")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Dry run mode")
	rootCmd.PersistentFlags().Bool("automatic", true, "Include automatic snapshots")
	rootCmd.PersistentFlags().Bool("manual", false, "Include manual snapshots")

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
