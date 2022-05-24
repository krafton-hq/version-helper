package version_object

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMerge(t *testing.T) {
	//language=json
	json1 := "{\"key\": \"value\", \"list\": [\"a\",{\"The end is\": \"nigh\"}]}"
	//language=json
	json2 := "{\"key\": \"override-value\",\"key2\": \"value2\", \"list\": [\"single-value\",{\"StructKey\": \"vALUE\"}]}"

	map1 := map[string]interface{}{}
	err := json.Unmarshal([]byte(json1), &map1)
	if err != nil {
		t.Fatal(err)
	}

	map2 := map[string]interface{}{}
	err = json.Unmarshal([]byte(json2), &map2)
	if err != nil {
		t.Fatal(err)
	}

	map3 := mergeMaps(map1, map2)
	fmt.Println(map3)
}
