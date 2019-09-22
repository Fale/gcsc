package clean

import "time"

type RetentionPolicy struct {
	Begin   time.Duration
	End     time.Duration
	Cadence time.Duration
}

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

type RetentionPolicies []RetentionPolicy

func (rs *RetentionPolicies) IsValid() bool {
	for _, r := range *rs {
		if !r.IsValid() {
			return false
		}
	}
	return true
}
