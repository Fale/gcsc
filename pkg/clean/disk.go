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

	"github.com/fale/gcsc/pkg/rp"
	compute "google.golang.org/api/compute/v1"
)

// Disks contains an array of disks to be processed
type Disks []Disk

// AddSnapshot allows to add a Snapshot to a Disks collection if the snapshot
// refers to a disk that is already present in the collection, a snapshot will
// be added to that disk. If the referred disk is not in the collection, it
// will be added.
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

// Disk is used to store basic information about a disk and its snapshots
type Disk struct {
	// Then name of the disk
	Name string
	// The GCP resource ID of the disk
	ID string
	// The list of snapshots of the disk
	Snapshots []compute.Snapshot
}

// Purgeable returns a list of snapshots that should be purged to complain with
// the passed retention policies
func (d *Disk) Purgeable(rps *rp.RetentionPolicies) []compute.Snapshot {
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
