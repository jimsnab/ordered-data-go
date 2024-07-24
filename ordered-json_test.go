package ordereddata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	_ "embed"
)

//go:embed "test_assets/ordered_json_1.json"
var orderedJson1Json []byte

//go:embed "test_assets/ordered_json_2.json"
var orderedJson2Json []byte

//go:embed "test_assets/ordered_json_nest_1.json"
var orderedJsonNest1Json []byte

//go:embed "test_assets/ordered_json_nest_2.json"
var orderedJsonNest2Json []byte

//go:embed "test_assets/ordered_json_nest_2.txt"
var orderedJsonNest2Txt string

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

func TestOrderedStringMapNesting(t *testing.T) {

	m := NewStringMap()
	err := json.Unmarshal(orderedJsonNest1Json, &m)
	if err != nil {
		t.Fatal(err)
	}

	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(raw, orderedJsonNest2Json) {
		t.Error("did not get expected output")
	}
}

func TestOrderedStringify(t *testing.T) {

	m := NewStringMap()
	err := json.Unmarshal(orderedJsonNest1Json, &m)
	if err != nil {
		t.Fatal(err)
	}

	s := m.String()

	if s != orderedJsonNest2Txt {
		t.Error("did not get expected output")
	}
}

func TestJsonEscape(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{`A string with "double quotes"`, `"A string with \"double quotes\""`},
		{`A string with a unicode character: \u2028`, `"A string with a unicode character: \\u2028"`},
		{`Plain string`, `"Plain string"`},
		{123, "123"},
		{true, "true"},
		{[]string{"a", "b", "c"}, `["a","b","c"]`},
	}

	for _, test := range tests {
		result := jsonEscape(test.input)
		if result != test.expected {
			t.Errorf("jsonEscape(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}
