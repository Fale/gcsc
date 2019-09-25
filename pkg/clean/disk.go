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
	"fmt"
	"sort"
	"time"

	compute "google.golang.org/api/compute/v1"
)

type Disks []Disk

func (ds *Disks) AddSnapshot(s compute.Snapshot) {
	for k, d := range *ds {
		if d.ID == s.SourceDiskId {
			(*ds)[k].Snapshots = append((*ds)[k].Snapshots, s)
			return
		}
	}
	*ds = append(*ds, Disk{
		Name:      s.SourceDisk,
		ID:        s.SourceDiskId,
		Snapshots: []compute.Snapshot{s},
	})
}

type Disk struct {
	Name      string
	ID        string
	Snapshots []compute.Snapshot
}

func (d *Disk) Purgeable(rps *RetentionPolicies) []compute.Snapshot {
	var ps []compute.Snapshot
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	for _, rp := range *rps {
		begin := now.Add(-1 * rp.Begin)
		end := now.Add(-1 * rp.End)
		for wb := begin; wb.After(end); wb = wb.Add(-1 * rp.Cadence) {
			drop := false
			sort.Slice(d.Snapshots, func(i, j int) bool {
				return d.Snapshots[i].CreationTimestamp < d.Snapshots[j].CreationTimestamp
			})
			for _, s := range d.Snapshots {
				ct, err := time.Parse("2006-01-02T15:04:05Z07:00", s.CreationTimestamp)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				if ct.Before(wb) && ct.After(wb.Add(-1*rp.Cadence)) {
					if drop {
						ps = append(ps, s)
					} else {
						drop = true
					}
				}
			}
		}
	}
	return ps
}
