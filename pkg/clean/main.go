package clean

import (
	"context"
	"errors"
	"fmt"

	compute "google.golang.org/api/compute/v1"
)

type Parameters struct {
	ProjectID          string
	CleanAutoBackups   bool
	CleanManualBackups bool
	DryRun             bool
	RetentionPolicies  RetentionPolicies
}

func Execute(p Parameters) error {
	if !p.RetentionPolicies.IsValid() {
		panic(errors.New("the retention policies are not valid"))
	}

	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		return err
	}

	var ds Disks
	err = computeService.Snapshots.List(p.ProjectID).Pages(ctx, func(ss *compute.SnapshotList) error {
		for _, s := range ss.Items {
			if s.AutoCreated && p.CleanAutoBackups {
				ds.AddSnapshot(*s)
			}
			if !s.AutoCreated && p.CleanManualBackups {
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
		pss = append(pss, d.Purgeable(&p.RetentionPolicies)...)
	}
	for _, ps := range pss {
		fmt.Printf("Deleting %s\n", ps.Name)
		if !p.DryRun {
			_, err = computeService.Snapshots.Delete(p.ProjectID, ps.Name).Do()
			if err != nil {
				fmt.Println(err)
				return nil
			}
		}
	}
	return nil
}
