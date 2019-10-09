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

package rp

import "time"

// RetentionPolicy is an object to describe a retention policy.
type RetentionPolicy struct {
	// Begin of the retention policy window, thinking from today backward
	Begin time.Duration
	// End of the retention policy window, thinking from today backward
	End time.Duration
	// Cadence of the expected backup sequence
	Cadence time.Duration
}

// IsValid validates that the RetentionPolicy is formally correct
func (r *RetentionPolicy) IsValid() bool {
	delta := r.End.Nanoseconds() - r.Begin.Nanoseconds()
	if delta < 0 {
		return false
	}
	if r.Cadence < 0 {
		return false
	}
	if r.Cadence != 0 {
		if div := delta % r.Cadence.Nanoseconds(); div != 0 {
			return false
		}
	}
	return true
}

// RetentionPolicies is used to aggregate multiple RetentionPolicy windows
type RetentionPolicies []RetentionPolicy

//IsValid validates that the RetentionPolicies set is formally correct
func (rs *RetentionPolicies) IsValid() bool {
	for _, r := range *rs {
		if !r.IsValid() {
			return false
		}
	}
	return true
}
