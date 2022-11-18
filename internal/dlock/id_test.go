package dlock

import "testing"

func TestDlockID_Key(t *testing.T) {
	type testcase struct {
		input DLockID
		want  string
	}

	tcs := map[string]testcase{
		"simple": {
			"table/default.tfstate",
			"default.tfstate",
		},
		"deep tfstate path": {
			"table/my/service/default.tfstate",
			"my/service/default.tfstate",
		},
		"empty": {
			"table",
			"",
		},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			tc := tc
			got := tc.input.Key()
			if got != tc.want {
				t.Fatalf("result mismatch, got:%s want:%s", got, tc.want)
			}
		})
	}
}

func TestDlockID_Table(t *testing.T) {
	type testcase struct {
		input DLockID
		want  string
	}

	tcs := map[string]testcase{
		"simple": {
			"table/default.tfstate",
			"table",
		},
		"deep tfstate path": {
			"table/my/service/default.tfstate",
			"table",
		},
		"empty": {
			"table",
			"table",
		},
	}

	for n, tc := range tcs {
		t.Run(n, func(t *testing.T) {
			tc := tc
			got := tc.input.Table()
			if got != tc.want {
				t.Fatalf("result mismatch, got:%s want:%s", got, tc.want)
			}
		})
	}
}
