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
