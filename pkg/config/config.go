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

type Config struct {
	ProjectID         string               `mapstructure:"project-id"`
	RetentionPolicies rp.RetentionPolicies `mapstructure:"retention-policies"`
	Automatic         bool                 `mapstructure:"automatic"`
	Manual            bool                 `mapstructure:"manual"`
	DryRun            bool                 `mapstructure:"dry-run"`
}

func Load(cmd *cobra.Command) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath("$HOME/.gcsc")

	v.AutomaticEnv()
	v.SetEnvPrefix("gcsc")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	v.BindPFlags(cmd.PersistentFlags())

	c := &Config{}
	if err := v.Unmarshal(c); err != nil {
		return nil, err
	}

	return c, nil
}
