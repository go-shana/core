package data

import (
	"testing"

	"github.com/huandu/go-assert"
)

func TestNameConverter(t *testing.T) {
	a := assert.New(t)
	cases := []struct {
		From          string
		To            string
		NameConverter NameConverter
	}{
		{ // Capitalize.
			"abc",
			"Abc",
			Capitalize,
		},
		{ // Capitalize empty string.
			"",
			"",
			Capitalize,
		},
		{ // Uncapitalize.
			"Abc",
			"abc",
			Uncapitalize,
		},
		{ // Uncapitalize empty string.
			"",
			"",
			Uncapitalize,
		},
		{ // SnakeCase.
			"AbcDef",
			"abc_def",
			SnakeCase,
		},
		{ // SnakeCase empty string.
			"",
			"",
			SnakeCase,
		},
		{ // CamelCase.
			"abc_def",
			"AbcDef",
			CamelCase,
		},
		{ // CamelCase empty string.
			"",
			"",
			CamelCase,
		},
	}

	for i, c := range cases {
		a.Use(&i, &c)
		a.Equal(c.To, c.NameConverter(c.From))
	}
}
