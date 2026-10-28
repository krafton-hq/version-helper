package version_object

import "testing"

func TestIsNewerRevision(t *testing.T) {
	cases := []struct {
		name               string
		current, uploading int
		want               bool
	}{
		{"latest has higher revision: keep it", 20, 10, true},
		{"uploading is newer: overwrite", 10, 20, false},
		{"same revision: overwrite", 15, 15, false},
		{"current has no revision: overwrite", 0, 10, false},
		{"uploading has no revision: overwrite", 20, 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isNewerRevision(c.current, c.uploading); got != c.want {
				t.Fatalf("isNewerRevision(%d, %d) = %v, want %v", c.current, c.uploading, got, c.want)
			}
		})
	}
}
