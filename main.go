package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	compute "google.golang.org/api/compute/v1"
)

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.gcsc")
	_ = viper.ReadInConfig()

	rootCmd := initCmd()

	if err := viper.BindPFlag("project-id", rootCmd.PersistentFlags().Lookup("project-id")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("automatic", rootCmd.PersistentFlags().Lookup("automatic")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("manual", rootCmd.PersistentFlags().Lookup("manual")); err != nil {
		panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func clean(cmd *cobra.Command, args []string) error {
	var rps RetentionPolicies
	err := viper.UnmarshalKey("retention-policies", &rps)
	if err != nil {
		panic(err)
	}
	if !rps.IsValid() {
		panic(errors.New("the retention policies are not valid"))
	}

	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return err
	}

	var ds Disks
	err = computeService.Snapshots.List(viper.GetString("project-id")).Pages(ctx, func(ss *compute.SnapshotList) error {
		for _, s := range ss.Items {
			if s.AutoCreated && viper.GetBool("automatic") {
				ds.AddSnapshot(*s)
			}
			if !s.AutoCreated && viper.GetBool("manual") {
				ds.AddSnapshot(*s)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	var pss []compute.Snapshot // nolint:prealloc
	for _, d := range ds {
		pss = append(pss, d.Purgeable(&rps)...)
	}
	for _, ps := range pss {
		fmt.Printf("Deleting %s\n", ps.Name)
		if !viper.GetBool("dry-run") {
			_, err = computeService.Snapshots.Delete(viper.GetString("project-id"), ps.Name).Do()
			if err != nil {
				fmt.Println(err)
				return nil
			}
		}
	}
	return nil
}
