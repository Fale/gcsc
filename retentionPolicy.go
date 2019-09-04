package main

import "time"

type RetentionPolicy struct {
	Begin   time.Duration
	End     time.Duration
	Cadence time.Duration
}

// TODO: Implementation
func (r *RetentionPolicy) IsValid() bool {
	return true
}

type RetentionPolicies []RetentionPolicy
