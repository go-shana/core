package data

import "strings"

// FieldTag is a parsed field tag.
// The format is:
//
//	data:"alias,opt1,opt2,..."
//
// Current supported options are:
//
//   - omitempty: ignore empty value.
//   - squash: squash this field.
//
// When alias is "-", current field will be skipped.
type FieldTag struct {
	Alias     string // The alias set in tag.
	Skipped   bool   // Skipped if the alias is "-".
	OmitEmpty bool   // Ignore empty value.
	Squash    bool   // Squash this field.
}

// ParseFieldTag parses alias and options from field tag.
// See doc in FieldTag for more details.
func ParseFieldTag(tag string) *FieldTag {
	opts := strings.Split(tag, ",")
	alias := strings.TrimSpace(opts[0])
	skipped := false
	omitEmpty := false
	squash := false

	for _, opt := range opts[1:] {
		switch opt {
		case "omitempty":
			omitEmpty = true
		case "squash":
			squash = true
		}
	}

	if alias == "-" {
		alias = ""
		skipped = true
	}

	return &FieldTag{
		Alias:     alias,
		Skipped:   skipped,
		OmitEmpty: omitEmpty,
		Squash:    squash,
	}
}
