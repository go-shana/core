package data

import (
	"testing"

	"github.com/huandu/go-assert"
)

func TestEncoder(t *testing.T) {
	cases := []struct {
		Value     any
		Data      Data
		OmitEmpty bool
	}{
		{ // Nil value.
			nil,
			Data{},
			false,
		},
		{ // Simple value.
			&AllValue{
				Int: 123,
			},
			Make(RawData{
				"int": 123,
			}),
			true,
		},
		{ // All data.
			allValues,
			fullData,
			false,
		},
	}
	a := assert.New(t)
	enc := &Encoder{
		TagName: "test",
	}

	for i, c := range cases {
		a.Use(&i, &c)

		enc.OmitEmpty = c.OmitEmpty
		a.Equal(c.Data, enc.Encode(c.Value))
	}
}
