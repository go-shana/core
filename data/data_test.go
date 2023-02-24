package data

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/huandu/go-assert"
)

func TestDataQuery(t *testing.T) {
	cases := []struct {
		Data   Data
		Fields []string
		Result any
	}{
		{ // Empty value.
			Data{},
			nil,
			nil,
		},
		{ // Simple case.
			Make(RawData{
				"a": 1,
			}),
			[]string{"a"},
			1,
		},
		{ // Field is not found.
			Make(RawData{
				"a": 1,
			}),
			[]string{"bad"},
			nil,
		},
		{ // Fake query.
			fullData,
			[]string{"fake", "query"},
			nil,
		},
		{ // Fake query.
			fullData,
			[]string{"fake.query"},
			float64(67.89),
		},
		{ // Complex query.
			fullData,
			[]string{"anonymous_type", "data_list", "2", "b"},
			true,
		},
		{ // Special map.
			Make(RawData{
				"m": map[int]any{
					123: map[uint]any{
						456: map[float64]any{
							0.5: "foo",
						},
					},
				},
			}),
			[]string{"m", "123", "456", "0.5"},
			"foo",
		},
	}

	a := assert.New(t)

loop:
	for i, c := range cases {
		a.Use(&i, &c)
		query := strings.Join(c.Fields, ".")

		a.Equal(c.Result, c.Data.Get(c.Fields...))

		for _, f := range c.Fields {
			if strings.Contains(f, ".") {
				continue loop
			}
		}

		a.Equal(c.Result, c.Data.Query(query))
	}
}

var (
	complexData = Make(RawData{
		"int":    123,
		"true":   true,
		"false":  false,
		"float":  12.34,
		"string": "string",
		"map": RawData{
			"m": "m",
		},
		"array": []RawData{
			{
				"d1": 1,
			},
			{
				"d2": "2",
			},
		},
		"ints":    []int{3, 2, 1},
		"floats":  []float64{5.5, 4.5, 3.5},
		"strings": []string{"s1", "s2", "s3"},
		"any":     []any{1, "2", 3.3},
	})
	complexDataJSON = `{
	"any": [
		1,
		"2",
		3.3
	],
	"array": [
		{
			"d1": 1
		},
		{
			"d2": "2"
		}
	],
	"false": false,
	"float": 12.34,
	"floats": [
		5.5,
		4.5,
		3.5
	],
	"int": 123,
	"ints": [
		3,
		2,
		1
	],
	"map": {
		"m": "m"
	},
	"string": "string",
	"strings": [
		"s1",
		"s2",
		"s3"
	],
	"true": true
}`
)

func TestDataString(t *testing.T) {
	cases := []struct {
		Data       Data
		PrettyJSON string
	}{
		{ // Empty value.
			Data{},
			"{}",
		},
		{ // Typical case.
			complexData,
			complexDataJSON,
		},
	}

	a := assert.New(t)

	for _, c := range cases {
		a.Use(&c)
		a.Equal(c.Data.JSON(true), c.PrettyJSON)

		buf := &bytes.Buffer{}
		a.NilError(json.Compact(buf, []byte(c.PrettyJSON)))
		str := buf.String()
		a.Equal(c.Data.String(), str)
	}
}

func TestDataJSONUnmarshal(t *testing.T) {
	cases := []struct {
		JSON     string
		Value    any
		HasError bool
	}{
		{ // Test against all kinds of data formats.
			`{"a":123, "data":{"int":123, "float":2.5, "strings":["s1", "s2"]}}`,
			&struct {
				A    int  `json:"a"`
				Data Data `json:"data"`
			}{
				A: 123,
				Data: Make(RawData{
					"int":     int64(123),
					"float":   float64(2.5),
					"strings": []any{"s1", "s2"},
				}),
			},
			false,
		},
		{ // Invalid value.
			`{"data":["s1", "s2"]}`,
			&struct {
				Data Data `json:"data"`
			}{},
			true,
		},
	}

	a := assert.New(t)

	for _, c := range cases {
		vt := reflect.ValueOf(c.Value).Type()
		actual := reflect.New(vt).Elem()
		err := json.Unmarshal([]byte(c.JSON), actual.Addr().Interface())

		if c.HasError {
			a.NonNilError(err)
		} else {
			a.NilError(err)
		}

		a.Equal(c.Value, actual.Interface())
	}
}

func TestParseQuery(t *testing.T) {
	a := assert.New(t)
	cases := []struct {
		Query  string
		Fields []string
	}{
		{ // Empty query.
			"",
			nil,
		},
		{ // Simple query.
			"a",
			[]string{"a"},
		},
		{ // Typical query.
			"a.b.c",
			[]string{"a", "b", "c"},
		},
		{ // Query with number.
			"a.1.c",
			[]string{"a", "1", "c"},
		},
		{ // Query with escape sequence.
			"a.b\\.c",
			[]string{"a", "b.c"},
		},
		{ // Query with backslash.
			"a.b\\\\.c",
			[]string{"a", "b\\", "c"},
		},
		{ // Query with escaped backslash.
			"a\\b.c",
			[]string{"a\\b", "c"},
		},
		{ // Query with dots.
			"...",
			[]string{"", "", "", ""},
		},
	}

	for i, c := range cases {
		a.Use(&i, &c)
		a.Equal(c.Fields, parseQuery(c.Query))
	}
}
