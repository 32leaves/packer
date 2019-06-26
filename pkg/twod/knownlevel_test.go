package twod_test

import (
	"github.com/32leaves/packer/pkg/twod"
	"github.com/go-test/deep"
	"testing"
)

func TestNextFitDecreasingHeight(t *testing.T) {
	tests := []struct {
		Name        string
		Problem     twod.Problem
		Solution    *twod.Solution
		ExpectedErr error
	}{
		{
			"unbound-height",
			twod.Problem{
				Sheet: twod.Sheet{Width: 1000},
				Items: []twod.Item{
					twod.Item{"a", 500, 200},
					twod.Item{"b", 495, 300},
					twod.Item{"c", 600, 400},
				},
			},
			&twod.Solution{
				Sheets: []twod.PackedSheet{
					twod.PackedSheet{
						Sheet: twod.Sheet{Width: 1000, Height: 0},
						Items: []twod.PackedItem{
							twod.PackedItem{twod.Item{"c", 600, 400}, 0, 0},
							twod.PackedItem{twod.Item{"b", 495, 300}, 0, 400},
							twod.PackedItem{twod.Item{"a", 500, 200}, 495, 400},
						},
					},
				},
				Cost: 700000,
			},
			nil,
		},
		{
			"bound-height",
			twod.Problem{
				Sheet: twod.Sheet{Width: 1000, Height: 900},
				Items: []twod.Item{
					twod.Item{"a", 500, 200},
					twod.Item{"b", 495, 300},
					twod.Item{"c", 600, 400},
				},
			},
			&twod.Solution{
				Sheets: []twod.PackedSheet{
					twod.PackedSheet{
						Sheet: twod.Sheet{Width: 1000, Height: 900},
						Items: []twod.PackedItem{
							twod.PackedItem{twod.Item{"c", 600, 400}, 0, 0},
							twod.PackedItem{twod.Item{"b", 495, 300}, 0, 400},
							twod.PackedItem{twod.Item{"a", 500, 200}, 495, 400},
						},
					},
				},
				Cost: 900000,
			},
			nil,
		},
	}

	for _, test := range tests {
		t.Log(test.Name)

		sol, err := twod.FirstFitDecreasingHeight(test.Problem)
		if err != nil {
			if err == test.ExpectedErr {
				continue
			}

			t.Errorf("[%s] unexpected error: %v", test.Name, err)
			continue
		}
		if test.ExpectedErr != nil {
			t.Errorf("[%s] expected error (%v) but got nil", test.Name, test.ExpectedErr)
			continue
		}

		diff := deep.Equal(sol, test.Solution)
		if diff != nil {
			t.Errorf("[%s] expected solution mismatch: %v", test.Name, diff)
		}
	}
}
