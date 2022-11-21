package dlock

import (
	"sync"
	"testing"
)

func TestDMutex_IsReadable(t *testing.T) {
	type testcase struct {
		arg DLockID

		// Input
		table DLockID
		mu    *dmutex

		want bool
	}

	tcs := map[string]testcase{
		"readable on first try": {
			DLockID("not-exist"),
			DLockID("table"),
			&dmutex{
				mu:        &sync.RWMutex{},
				isReading: true,
				isWriting: true,
			},
			true,
		},
		"not readable if table is writing": {
			DLockID("table"),
			DLockID("table"),
			&dmutex{
				mu:        &sync.RWMutex{},
				isWriting: true,
			},
			false,
		},
		"not readable if table is reading": {
			DLockID("table"),
			DLockID("table"),
			&dmutex{
				mu:        &sync.RWMutex{},
				isReading: true,
				isWriting: true,
			},
			false,
		},
		"readable if table is not reading": {
			DLockID("table"),
			DLockID("table"),
			&dmutex{
				mu:        &sync.RWMutex{},
				isReading: false,
			},
			true,
		},
		"readable if table is not writing": {
			DLockID("table"),
			DLockID("table"),
			&dmutex{
				mu:        &sync.RWMutex{},
				isWriting: false,
			},
			true,
		},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			tc := tc

			dmu := &DMutex{
				mu: &sync.RWMutex{},
				mus: map[DLockID]*dmutex{
					tc.table: tc.mu,
				},
			}

			got := dmu.IsReadable(tc.arg)
			if got != tc.want {
				t.Fatalf("result mismatch, got:%t want:%t", got, tc.want)
			}
		})
	}
}

func TestDMutex_IsWritable(t *testing.T) {
	type testcase struct {
		arg DLockID

		// Input
		table DLockID
		mu    *dmutex

		want bool
	}

	tcs := map[string]testcase{
		"writable on first try": {
			DLockID("not-exist"),
			DLockID("table"),
			&dmutex{
				mu:        &sync.RWMutex{},
				isWriting: true,
			},
			true,
		},
		"not writable if table is writing": {
			DLockID("table"),
			DLockID("table"),
			&dmutex{
				mu:        &sync.RWMutex{},
				isWriting: true,
			},
			false,
		},
		"writable if table is not writing": {
			DLockID("table"),
			DLockID("table"),
			&dmutex{
				mu:        &sync.RWMutex{},
				isWriting: false,
			},
			true,
		},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			tc := tc

			dmu := &DMutex{
				mu: &sync.RWMutex{},
				mus: map[DLockID]*dmutex{
					tc.table: tc.mu,
				},
			}

			got := dmu.IsWritable(tc.arg)
			if got != tc.want {
				t.Fatalf("result mismatch, got:%t want:%t", got, tc.want)
			}
		})
	}
}
