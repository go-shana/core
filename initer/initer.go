// Package init provides a function to initialize data and its struct fields in depth-first order.
package initer

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-shana/core/errors"
)

var errInit = errors.New("fail to initialize data")

// Initer represents a data which can be automatically initialized by Shana.
type Initer interface {
	Init(ctx context.Context)
}

type initerVisitMap map[reflect.Type]struct{}

var typeOfIniter = reflect.TypeOf((*Initer)(nil)).Elem()

// Init initializes data and its struct fields in depth-first order.
//
// Init is aware of several struct field tags as stated following.
//
//   - `init:"disabled"`: Fully disables the initialization on the field.
//   - `init:"omitempty"`: Skip initialization if the field is an empty pointer.
func Init(ctx context.Context, data any) (err error) {
	if data == nil {
		err = errInit
		return
	}

	val := reflect.ValueOf(data)

	if kind := val.Kind(); (kind == reflect.Pointer || kind == reflect.Interface) && val.IsNil() {
		err = errInit
		return
	}

	err = initValue(ctx, val)
	return
}

// initValue initializes val if val and/or val's struct fields implement Initer.
// If a nil value implements Initer, the value will be set to a zero new value before calling Init.
func initValue(ctx context.Context, val reflect.Value) (err error) {
	defer errors.Handle(&err)
	var init reflect.Value

	allocateValue(val)

	// If pointer to val type implements Initer, use it.
	if val.CanAddr() {
		ptr := val.Addr()

		if t := ptr.Type(); t.Implements(typeOfIniter) {
			init = ptr
		}
	}

	elem := val

	// Allocate new value if val implements Initer and is nil.
	if !init.IsValid() {
		for elem.IsValid() {
			if t := elem.Type(); t.Implements(typeOfIniter) {
				init = elem
				break
			}

			if elem.Kind() != reflect.Pointer {
				break
			}

			elem = elem.Elem()
			allocateValue(elem)
		}
	}

	for elem.IsValid() && elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
		allocateValue(elem)
	}

	if elem.Kind() == reflect.Struct {
		t := elem.Type()
		num := elem.NumField()

		// Depth-first initialization.
		for i := 0; i < num; i++ {
			stField := t.Field(i)

			if !stField.IsExported() {
				continue
			}

			if isStructFieldDisabled(&stField) {
				continue
			}

			field := elem.Field(i)

			if isStructFieldOmitempty(&stField) {
				if field.Kind() == reflect.Pointer && field.IsNil() {
					continue
				}
			}

			errors.Check(initValue(ctx, field))
		}
	}

	if init.IsValid() {
		init.Interface().(Initer).Init(ctx)
	}

	return
}

func allocateValue(val reflect.Value) {
	if !val.CanSet() {
		return
	}

	t := val.Type()

	switch val.Kind() {
	case reflect.Pointer:
		if !val.IsNil() || !shouldAllocate(t) {
			break
		}

		val.Set(reflect.New(t.Elem()))

	case reflect.Map:
		if !val.IsNil() || !shouldAllocate(t) {
			break
		}

		val.Set(reflect.MakeMap(t))
	}
}

func isStructFieldDisabled(field *reflect.StructField) bool {
	tags := strings.Split(field.Tag.Get("init"), ",")

	for _, tag := range tags {
		if tag == "disabled" {
			return true
		}
	}

	return false
}

func isStructFieldOmitempty(field *reflect.StructField) bool {
	tags := strings.Split(field.Tag.Get("init"), ",")

	for _, tag := range tags {
		if tag == "omitempty" {
			return true
		}
	}

	return false
}

// shouldAllocate checks whether t or any of t's fields implements Initer.
//
// Note: shana will compute the result at compile time.
// There is no need to build a global cache for the result.
func shouldAllocate(t reflect.Type) bool {
	visited := initerVisitMap{}
	return shouldAllocateRecursive(t, visited)
}

func shouldAllocateRecursive(t reflect.Type, visited initerVisitMap) bool {
	for {
		if t.Implements(typeOfIniter) {
			return true
		}

		if _, ok := visited[t]; ok {
			return false
		}

		visited[t] = struct{}{}

		if t.Kind() != reflect.Pointer {
			break
		}

		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return false
	}

	num := t.NumField()
	hasIniter := false

	for i := 0; i < num; i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		if isStructFieldDisabled(&field) {
			continue
		}

		if _, ok := visited[field.Type]; ok {
			continue
		}

		if field.Type.Implements(typeOfIniter) {
			hasIniter = true
			continue
		}

		if reflect.PtrTo(field.Type).Implements(typeOfIniter) {
			hasIniter = true
			continue
		}

		visited[field.Type] = struct{}{}
	}

	return hasIniter
}
