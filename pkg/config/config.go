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

package config

import (
	"strings"

	"github.com/fale/gcsc/pkg/rp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config is used to provide the configuration to GCSC application
type Config struct {
	// GCP Project ID to operate on
	ProjectID string `mapstructure:"project-id"`
	// A set of RetentionPolicy to apply to the snapshots
	RetentionPolicies rp.RetentionPolicies `mapstructure:"retention-policies"`
	// Include automatically generated snapshot in the calculations and pruning
	Automatic bool `mapstructure:"automatic"`
	// Include manually generated snapshots in the calculations and pruning
	Manual bool `mapstructure:"manual"`
	// Do not perform any real operations
	DryRun bool `mapstructure:"dry-run"`
}

// Load returns a Config item with the values gathered from the configurations
// file, or the CLI or the ENV VARS
func Load(cmd *cobra.Command) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("$HOME/.gcsc")

	v.AutomaticEnv()
	v.SetEnvPrefix("gcsc")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := v.BindPFlags(cmd.PersistentFlags()); err != nil {
		return nil, err
	}

	c := &Config{}
	if err := v.Unmarshal(c); err != nil {
		return nil, err
	}

	return c, nil
}
