package clean_test

import (
	"testing"
	"time"

	"github.com/fale/gcsc/pkg/clean"
)

func TestDataBeginTime(t *testing.T) {
	tests := []struct {
		rp    clean.RetentionPolicy
		valid bool
	}{
		{ // All zeros
			rp: clean.RetentionPolicy{
				Begin:   0,
				End:     0,
				Cadence: 0,
			},
			valid: true,
		},
		{ // 1h slot, 1h cadence
			rp: clean.RetentionPolicy{
				Begin:   0,
				End:     time.Hour,
				Cadence: time.Hour,
			},
			valid: true,
		},
		{ // 1h slot, 1h cadence
			rp: clean.RetentionPolicy{
				Begin:   time.Hour,
				End:     2 * time.Hour,
				Cadence: time.Hour,
			},
			valid: true,
		},
		{ // 3h slot, 2h cadence
			rp: clean.RetentionPolicy{
				Begin:   0,
				End:     3 * time.Hour,
				Cadence: 2 * time.Hour,
			},
			valid: false,
		},
		{ // negative cadence
			rp: clean.RetentionPolicy{
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
