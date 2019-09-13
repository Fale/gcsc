package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func config(rootCmd *cobra.Command) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.gcsc")
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
