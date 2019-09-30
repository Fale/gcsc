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

package rp_test

import (
	"testing"
	"time"

	"github.com/fale/gcsc/pkg/rp"
)

func TestDataBeginTime(t *testing.T) {
	tests := []struct {
		rp    rp.RetentionPolicy
		valid bool
	}{
		{ // All zeros
			rp: rp.RetentionPolicy{
				Begin:   0,
				End:     0,
				Cadence: 0,
			},
			valid: true,
		},
		{ // 1h slot, 1h cadence
			rp: rp.RetentionPolicy{
				Begin:   0,
				End:     time.Hour,
				Cadence: time.Hour,
			},
			valid: true,
		},
		{ // 1h slot, 1h cadence
			rp: rp.RetentionPolicy{
				Begin:   time.Hour,
				End:     2 * time.Hour,
				Cadence: time.Hour,
			},
			valid: true,
		},
		{ // 3h slot, 2h cadence
			rp: rp.RetentionPolicy{
				Begin:   0,
				End:     3 * time.Hour,
				Cadence: 2 * time.Hour,
			},
			valid: false,
		},
		{ // negative cadence
			rp: rp.RetentionPolicy{
				Begin:   0,
				End:     time.Hour,
				Cadence: -1 * time.Hour,
			},
			valid: false,
		},
	}
	for _, test := range tests {
		if test.rp.IsValid() != test.valid {
			t.Errorf("expecting %v, received %v", test.valid, test.rp.IsValid())
		}
	}
}
