package api

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kerraform/kerranamodb/internal/id"
)

func TestPutInput_GetLockID(t *testing.T) {
	type testcase struct {
		item map[string]map[string]string
		want id.LockID
		err  bool
	}

	tcs := map[string]*testcase{
		"ok": {
			item: map[string]map[string]string{
				"LockID": {
					"S": "kerranamodb/tfstate.tfstate",
				},
			},
			want: id.LockID("kerranamodb/tfstate.tfstate"),
		},
		"missing S": {
			item: map[string]map[string]string{
				"LockID": {},
			},
			err: true,
		},
		"missing lock ID": {
			item: map[string]map[string]string{},
			err:  true,
		},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			tc := tc
			t.Parallel()

			i := &PutInput{
				Item: tc.item,
			}

			got, err := i.GetLockID()
			if err != nil {
				if !tc.err {
					t.Fatalf("unexpected error: %v", err)
				}

				return
			}

			if tc.err {
				t.Fatalf("test should fail")
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("unexpected result, diff(+got,-want): %s\n", diff)
			}
		})
	}
}
