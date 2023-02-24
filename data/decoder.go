package data

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/go-shana/core/errors"
	"github.com/huandu/go-clone"
)

var errDecodeInvalidValue = errors.New("cannot decode to an invalid value")

// Decoder decodes a Data and updates the value of target struct in depth.
type Decoder struct {
	TagName       string        // The field tag used in decoder. The default value is "data".
	NameConverter NameConverter // The function used to convert field name to another name. The default value is nil.
}

// Decode decodes d and updates v in depth.
func (dec *Decoder) Decode(d Data, v any) error {
	from := reflect.ValueOf(d.data)
	to := reflect.ValueOf(v)
	return dec.decode(from, to)
}

// DecodeQuery decodes the part of d matching the query and updates v in depth.
// See `Data#Query` for the syntax of the query.
func (dec *Decoder) DecodeQuery(d Data, query string, v any) error {
	from := reflect.ValueOf(d.Query(query))
	to := reflect.ValueOf(v)
	return dec.decode(from, to)
}

// DecodeField decodes the part of d matching the fields and updates v in depth.
// See `Data#Get` for the rule using the fields.
func (dec *Decoder) DecodeField(d Data, fields []string, v any) error {
	from := reflect.ValueOf(d.Get(fields...))
	to := reflect.ValueOf(v)
	return dec.decode(from, to)
}

// decode decodes the value of from and updates the to.
// The to.CanSet() must be true.
//
// As decode is a private function, we can assume the value of from must be Data or a parsed Data.
// So the from cannot be and cannot include any of struct, chan, func and ptr.
func (dec *Decoder) decode(from reflect.Value, to reflect.Value) error {
	if to.Kind() == reflect.Pointer {
		to = to.Elem()
	}

	if !to.IsValid() {
		return errDecodeInvalidValue
	}

	if !to.CanSet() {
		return fmt.Errorf("cannot decode to a value of type %v which is not settable", to.Type())
	}

	switch from.Kind() {
	case reflect.Invalid:
		// if from == nil, skip the decoding process.
		return nil
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Pointer, reflect.Slice, reflect.Map:
		if from.IsNil() {
			return nil
		}
	}

	for to.Kind() == reflect.Pointer {
		if to.IsNil() {
			to.Set(reflect.New(to.Type().Elem()))
		}

		to = to.Elem()
	}

	for from.Kind() == reflect.Interface {
		from = from.Elem()
	}

	// Decodes well-known types.
	switch to.Type() {
	case typeOfDuration:
		if from.Kind() != reflect.String {
			return fmt.Errorf("cannot decode a value of type %v from %v", to.Type(), from.Type())
		}

		if str := from.String(); str == "" {
			to.SetInt(0)
		} else {
			dur, err := time.ParseDuration(from.String())

			if err != nil {
				return err
			}

			to.SetInt(int64(dur))
		}

		return nil

	case typeOfTime:
		if from.Type() != typeOfTime {
			return fmt.Errorf("cannot decode a value of type %v from %v", to.Type(), from.Type())
		}

		to.Set(from)
		return nil
	}

	// Decodes primitive types.
	switch to.Kind() {
	case reflect.Bool:
		switch from.Kind() {
		case reflect.Bool:
			to.SetBool(from.Bool())
			return nil
		}

	case reflect.String:
		switch from.Kind() {
		case reflect.String:
			to.SetString(from.String())
			return nil

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i := from.Int()
			to.SetString(strconv.FormatInt(i, 10))
			return nil

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ui := from.Uint()
			to.SetString(strconv.FormatUint(ui, 10))
			return nil

		case reflect.Float32, reflect.Float64:
			f := from.Float()
			to.SetString(strconv.FormatFloat(f, 'f', -1, 64))
			return nil

		case reflect.Bool:
			to.SetString(strconv.FormatBool(from.Bool()))
			return nil
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch from.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i := from.Int()

			if to.OverflowInt(i) {
				return fmt.Errorf("cannot decode value of type %v from %v due to overflow", to.Type(), i)
			}

			to.SetInt(i)
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ui := from.Uint()
			i := int64(ui)

			if ui > math.MaxInt64 || to.OverflowInt(i) {
				return fmt.Errorf("cannot decode value of type %v from %v due to overflow", to.Type(), ui)
			}

			to.SetInt(i)
			return nil
		case reflect.Float32, reflect.Float64:
			f := from.Float()
			i := int64(f)

			if f != math.Round(f) {
				return fmt.Errorf("cannot decode value of type %v from a float number %v", to.Type(), f)
			}

			if f > math.MaxInt64 || to.OverflowInt(i) {
				return fmt.Errorf("cannot decode value of type %v from %v due to overflow", to.Type(), f)
			}

			to.SetInt(i)
			return nil
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch from.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i := from.Int()
			ui := uint64(i)

			if i < 0 || to.OverflowUint(ui) {
				return fmt.Errorf("cannot decode value of type %v from %v due to overflow", to.Type(), i)
			}

			to.SetUint(ui)
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ui := from.Uint()

			if to.OverflowUint(ui) {
				return fmt.Errorf("cannot decode value of type %v from %v due to overflow", to.Type(), ui)
			}

			to.SetUint(ui)
			return nil
		case reflect.Float32, reflect.Float64:
			f := from.Float()
			ui := uint64(f)

			if f != math.Round(f) {
				return fmt.Errorf("cannot decode value of type %v from a float number %v", to.Type(), f)
			}

			if f < 0 || f > math.MaxUint64 || to.OverflowUint(ui) {
				return fmt.Errorf("cannot decode value of type %v from %v due to overflow", to.Type(), f)
			}

			to.SetUint(ui)
			return nil
		}

	case reflect.Float32, reflect.Float64:
		switch from.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i := from.Int()
			f := float64(i)
			to.SetFloat(f)
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ui := from.Uint()
			f := float64(ui)
			to.SetFloat(f)
			return nil
		case reflect.Float32, reflect.Float64:
			f := from.Float()

			if to.OverflowFloat(f) {
				return fmt.Errorf("cannot decode value of type %v from %v due to overflow", to.Type(), f)
			}

			to.SetFloat(f)
			return nil
		}

	case reflect.Complex64, reflect.Complex128:
		switch from.Kind() {
		case reflect.Complex64, reflect.Complex128:
			cmplx := from.Complex()

			if to.OverflowComplex(cmplx) {
				return fmt.Errorf("cannot decode value of type %v from %v due to overflow", to.Type(), cmplx)
			}

			to.SetComplex(cmplx)
			return nil
		}

	case reflect.Array:
		switch from.Kind() {
		case reflect.Array, reflect.Slice:
			fromLen := from.Len()
			toLen := to.Len()

			if fromLen > toLen {
				return fmt.Errorf("cannot decode value of type %v due to no enough room to store %v element(s)", to.Type(), fromLen)
			}

			for i := 0; i < fromLen; i++ {
				v := to.Index(i)

				if err := dec.decode(from.Index(i), v); err != nil {
					return err
				}
			}

			return nil
		}

	case reflect.Slice:
		switch from.Kind() {
		case reflect.Array, reflect.Slice:
			fromLen := from.Len()
			toType := to.Type()
			val := reflect.MakeSlice(toType, fromLen, fromLen)

			for i := 0; i < fromLen; i++ {
				v := val.Index(i)

				if err := dec.decode(from.Index(i), v); err != nil {
					return err
				}
			}

			to.Set(val)
			return nil
		}

	case reflect.Map:
		switch from.Kind() {
		case reflect.Map:
			toType := to.Type()
			toKeyType := toType.Key()
			toElemType := toType.Elem()

			if toKeyType.Kind() != reflect.String {
				return fmt.Errorf("cannot decode a value of type %v whose key is not string", to.Type())
			}

			val := reflect.MakeMap(toType)
			iter := from.MapRange()

			for iter.Next() {
				v := reflect.New(toElemType).Elem()

				if err := dec.decode(iter.Value(), v.Addr()); err != nil {
					return err
				}

				val.SetMapIndex(iter.Key(), v)
			}

			to.Set(val)
			return nil
		}

	case reflect.Struct:
		if to.Type().AssignableTo(typeOfData) {
			d := Data{}

			if err := dec.decode(from, reflect.ValueOf(&d.data)); err != nil {
				return err
			}

			if d.Len() == 0 {
				d = emptyData
			}

			to.Set(reflect.ValueOf(d))
			return nil
		}

		switch from.Kind() {
		case reflect.Map:
			numField := to.NumField()
			toType := to.Type()

			for i := 0; i < numField; i++ {
				f := toType.Field(i)
				fv := to.Field(i)

				if !fv.CanSet() || !fv.CanAddr() {
					continue
				}

				tagName := dec.TagName

				if tagName == "" {
					tagName = defaultTagName
				}

				tag := f.Tag.Get(tagName)
				ft := ParseFieldTag(tag)

				if ft.Skipped {
					continue
				}

				// If it's necessary to squash fields and the type of the ft is a struct or pointer to struct,
				// use from to update the fv.
				if ft.Squash {
					fieldType := f.Type

					for fieldType.Kind() == reflect.Pointer {
						fieldType = fieldType.Elem()
					}

					if fieldType.Kind() == reflect.Struct {
						if err := dec.decode(from, fv.Addr()); err != nil {
							return err
						}

						continue
					}
				}

				k := f.Name

				if ft.Alias == "" {
					if dec.NameConverter != nil {
						k = dec.NameConverter(k)
					}
				} else {
					k = ft.Alias
				}

				kv := from.MapIndex(reflect.ValueOf(k))

				if !kv.IsValid() {
					continue
				}

				if err := dec.decode(kv, fv.Addr()); err != nil {
					return err
				}
			}

			return nil
		}

	case reflect.Interface:
		fromType := from.Type()
		toType := to.Type()

		if !fromType.Implements(toType) {
			return fmt.Errorf("cannot decode an interface value of type %v from %v", toType, fromType)
		}

		to.Set(reflect.ValueOf(clone.Clone(from.Interface())))
		return nil
	}

	return fmt.Errorf("cannot decode a value of type %v from %v", to.Type(), from.Type())
}
