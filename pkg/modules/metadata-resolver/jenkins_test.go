package metadata_resolver

import (
	"fmt"
	"testing"
)

func TestJenkinsResolver_ResolveBuildMetadata(t *testing.T) {
	jenkins := &JenkinsResolver{}

	meta, err := jenkins.ResolveBuildMetadata()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", meta)
}
