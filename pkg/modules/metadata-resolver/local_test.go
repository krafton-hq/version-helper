package metadata_resolver

import (
	"fmt"
	"testing"

	log_helper "github.com/krafton-hq/version-helper/pkg/log-helper"
)

func TestLocalResolver_ResolveBuildMetadata(t *testing.T) {
	log_helper.Initialize()

	local := &LocalResolver{}

	meta, err := local.ResolveBuildMetadata()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", meta)
}
