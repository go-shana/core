package data

import (
	"strings"

	"github.com/huandu/xstrings"
)

// NameConverter is a function that converts field name to another name.
type NameConverter func(name string) string

// Capitalize converts field name to capitalized name.
//
// For example, "fieldName" will be converted to "FieldName".
func Capitalize(name string) string {
	if len(name) == 0 {
		return name
	}

	return strings.ToUpper(name[:1]) + name[1:]
}

// Uncapitalize converts field name to uncapitalized name.
//
// For example, "FieldName" will be converted to "fieldName".
func Uncapitalize(name string) string {
	if len(name) == 0 {
		return name
	}

	return strings.ToLower(name[:1]) + name[1:]
}

// SnakeCase converts field name from camel case to snake case.
//
// For example, "FieldName" will be converted to "field_name".
func SnakeCase(name string) string {
	return xstrings.ToSnakeCase(name)
}

// CamelCase converts field name from snake case to camel case.
//
// For example, "field_name" will be converted to "fieldName".
func CamelCase(name string) string {
	return xstrings.ToCamelCase(name)
}
