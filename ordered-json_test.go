package ordereddata

import (
	"bytes"
	"encoding/json"
	"testing"

	_ "embed"
)

//go:embed "test_assets/ordered_json_1.json"
var orderedJson1Json []byte

//go:embed "test_assets/ordered_json_2.json"
var orderedJson2Json []byte

func TestOrderedJson(t *testing.T) {

	m := NewOrderedMap[string, any]()
	err := json.Unmarshal(orderedJson1Json, &m)
	if err != nil {
		t.Fatal(err)
	}

	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(raw, orderedJson2Json) {
		t.Error("did not get expected output")
	}
}

func TestOrderedStringMap(t *testing.T) {

	m := NewStringMap()
	err := json.Unmarshal(orderedJson1Json, &m)
	if err != nil {
		t.Fatal(err)
	}

	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(raw, orderedJson2Json) {
		t.Error("did not get expected output")
	}
}
