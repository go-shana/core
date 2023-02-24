package data

import (
	"testing"

	"github.com/huandu/go-assert"
)

func TestParseFieldTag(t *testing.T) {
	cases := []struct {
		Tag      string
		FieldTag *FieldTag
	}{
		{ // Empty tag.
			"",
			&FieldTag{},
		},
		{ // Only alias.
			"abc",
			&FieldTag{
				Alias: "abc",
			},
		},
		{ // Pure options.
			",omitempty,not-valid,squash,",
			&FieldTag{
				OmitEmpty: true,
				Squash:    true,
			},
		},
		{ // Ignore "-"".
			"-",
			&FieldTag{
				Skipped: true,
			},
		},
		{ // All options.
			"a1_b2,squash,omitempty",
			&FieldTag{
				Alias:     "a1_b2",
				OmitEmpty: true,
				Squash:    true,
			},
		},
	}

	for _, c := range cases {
		expected := c.FieldTag
		actual := ParseFieldTag(c.Tag)
		assert.AssertEqual(t, expected, actual)
	}
}
