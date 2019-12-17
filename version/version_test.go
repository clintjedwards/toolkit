package version

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	testTable := []struct {
		name           string
		input          string
		expectedOutput Info
	}{
		{
			"properly formatted",
			"v0.0.1_1553466344_55eaf31",
			Info{
				Semver: "v0.0.1",
				Epoch:  "1553466344",
				Hash:   "55eaf31",
			},
		},
		{
			"missing a section",
			"v0.0.1_1553466344_",
			Info{
				Semver: "v0.0.1",
				Epoch:  "1553466344",
				Hash:   "",
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			output, err := Parse(tc.input)
			if err != nil {
				t.Errorf("Incorrect output from parse: %s", err)
			}
			if !cmp.Equal(tc.expectedOutput, output) {
				t.Errorf("Incorrect output from parse; Diff:\n %s", cmp.Diff(tc.expectedOutput, output))
			}

		})
	}

}
