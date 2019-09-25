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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func config(rootCmd *cobra.Command) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.gcsc")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("gcsc")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	_ = viper.ReadInConfig()

	if err := viper.BindPFlag("project-id", rootCmd.PersistentFlags().Lookup("project-id")); err != nil {
		return err
	}
	if err := viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run")); err != nil {
		return err
	}
	if err := viper.BindPFlag("automatic", rootCmd.PersistentFlags().Lookup("automatic")); err != nil {
		return err
	}
	if err := viper.BindPFlag("manual", rootCmd.PersistentFlags().Lookup("manual")); err != nil {
		return err
	}
	return nil
}
