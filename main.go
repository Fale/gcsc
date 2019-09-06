package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	compute "google.golang.org/api/compute/v1"
)

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.snapshot-cleaner")
	_ = viper.ReadInConfig()

	rootCmd.PersistentFlags().StringP("project-id", "p", "", "Google Cloud Project ID")
	if err := viper.BindPFlag("project-id", rootCmd.PersistentFlags().Lookup("project-id")); err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().Bool("dry-run", false, "Dry run mode")
	if err := viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run")); err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().Bool("automatic", true, "Include automatic backups")
	if err := viper.BindPFlag("automatic", rootCmd.PersistentFlags().Lookup("automatic")); err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().Bool("manual", false, "Include manual backups")
	if err := viper.BindPFlag("manual", rootCmd.PersistentFlags().Lookup("manual")); err != nil {
		panic(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:  "snapshot-cleaner",
	RunE: clean,
}

func clean(cmd *cobra.Command, args []string) error {
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

	// TODO: Make those configurable from config file
	rps := RetentionPolicies{
		RetentionPolicy{
			Begin:   0 * 24 * time.Hour,
			End:     7 * 24 * time.Hour,
			Cadence: time.Hour,
		},
		RetentionPolicy{
			Begin:   7 * 24 * time.Hour,
			End:     14 * 24 * time.Hour,
			Cadence: 24 * time.Hour,
		},
		RetentionPolicy{
			Begin:   14 * 24 * time.Hour,
			End:     63 * 24 * time.Hour,
			Cadence: 7 * 24 * time.Hour,
		},
		RetentionPolicy{
			Begin:   63 * 24 * time.Hour,
			End:     1063 * 24 * time.Hour,
			Cadence: 1000 * 24 * time.Hour,
		},
	}
	var pss []compute.Snapshot
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
