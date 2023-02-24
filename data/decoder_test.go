package data

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/huandu/go-assert"
)

func TestDecode(t *testing.T) {
	cases := []struct {
		Fields   []string
		Data     Data
		AllValue any
		Value    any
	}{
		{ // Nil value.
			nil,
			emptyData,
			(*AllValue)(nil),
			(*AllValue)(nil),
		},
		{ // Simple value.
			nil,
			Make(RawData{
				"int": 123,
			}),
			&AllValue{
				Int:        123,
				SquashType: &SquashType{},
			},
			&AllValue{
				Int:        123,
				SquashType: &SquashType{},
			},
		},
		{ // Numeric value.
			nil,
			Make(RawData{
				"int":      uint(345),
				"uint32":   int64(567),
				"uint64":   123456.,
				"duration": "13m20s",
			}),
			&AllValue{
				Int:      345,
				Duration: 13*time.Minute + 20*time.Second,
				SquashType: &SquashType{
					Uint32: 567,
					Uint64: 123456,
				},
			},
			&AllValue{
				Int:      345,
				Duration: 13*time.Minute + 20*time.Second,
				SquashType: &SquashType{
					Uint32: 567,
					Uint64: 123456,
				},
			},
		},
		{ // Full data test.
			nil,
			fullData,
			allValues,
			allValues,
		},
		{ // Query test.
			[]string{"sub_type", "int64"},
			fullData,
			allValues,
			int64(-64),
		},
		{ // Invalid query.
			[]string{"fake", "query"},
			fullData,
			allValues,
			float64(0),
		},
	}
	a := assert.New(t)
	dec := &Decoder{
		TagName: "test",
	}

	for i, c := range cases {
		a.Use(&i, &c)

		query := strings.Join(c.Fields, ".")
		avt := reflect.ValueOf(c.AllValue).Type()
		vt := reflect.ValueOf(c.Value).Type()

		{
			expected := c.AllValue
			actual := reflect.New(avt).Elem()
			a.NilError(dec.Decode(c.Data, actual.Addr().Interface()))
			a.Equal(expected, actual.Interface())
		}

		{
			expected := c.Value
			actual := reflect.New(vt).Elem()
			a.NilError(dec.DecodeField(c.Data, c.Fields, actual.Addr().Interface()))
			a.Equal(expected, actual.Interface())
		}

		{
			expected := c.Value
			actual := reflect.New(vt).Elem()
			a.NilError(dec.DecodeQuery(c.Data, query, actual.Addr().Interface()))
			a.Equal(expected, actual.Interface())
		}
	}
}
