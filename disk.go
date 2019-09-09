package main

import (
	"fmt"
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
