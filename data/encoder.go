package data

import (
	"encoding/json"
	"reflect"

	"github.com/huandu/go-clone"
)

// Encoder encodes any data to Data
type Encoder struct {
	TagName       string        // The field tag used in encoder. The default value is "data".
	OmitEmpty     bool          // If true, empty value will be omitted. The default value is false.
	NameConverter NameConverter // The function used to convert field name to another name. The default value is nil.
}

// Encode encodes any data to Data.
//
// Only the following types can be encoded successfully. If v is not one of them, Encode may return nil.
//   - Go struct and struct pointer.
//   - Any map[string]T, T can be any type.
func (enc *Encoder) Encode(v any) Data {
	if v == nil {
		return emptyData
	}

	val := reflect.ValueOf(v)
	t := val.Type()

	for kind := t.Kind(); kind == reflect.Pointer || kind == reflect.Interface; kind = t.Kind() {
		t = t.Elem()
		val = val.Elem()
	}

	d := enc.encodeValue(val)
	return Data{
		data: d,
	}
}

func (enc *Encoder) encodeValue(val reflect.Value) RawData {
	switch val.Kind() {
	case reflect.Map:
		return enc.encodeMap(val)
	case reflect.Struct:
		return enc.encodeStruct(val)
	}

	return nil
}

func (enc *Encoder) encodeMap(val reflect.Value) RawData {
	t := val.Type()

	if t.Key().Kind() != reflect.String {
		return nil
	}

	d := RawData{}

	if val.Len() == 0 {
		return nil
	}

	iter := val.MapRange()

	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		encoded := enc.encodeMapValue(v)

		if encoded.IsValid() {
			d[k.String()] = encoded.Interface()
		} else {
			d[k.String()] = nil
		}
	}

	return d
}

func (enc *Encoder) encodeStruct(val reflect.Value) RawData {
	d := RawData{}
	enc.encodeStructToData(val, d)
	return d
}

func (enc *Encoder) encodeStructToData(val reflect.Value, d RawData) {
	if val.Type().AssignableTo(typeOfData) {
		merge(reflect.ValueOf(d), val.Convert(typeOfData).Interface().(Data).data)
		return
	}

	t := val.Type()
	l := t.NumField()

	for i := 0; i < l; i++ {
		f := t.Field(i)
		tagName := enc.TagName

		if tagName == "" {
			tagName = defaultTagName
		}

		tag := f.Tag.Get(tagName)
		ft := ParseFieldTag(tag)

		if ft.Skipped {
			continue
		}

		k := f.Name

		if ft.Alias == "" {
			if enc.NameConverter != nil {
				k = enc.NameConverter(k)
			}
		} else {
			k = ft.Alias
		}

		fv := val.Field(i)
		encoded := enc.encodeMapValue(fv)

		if (ft.OmitEmpty || enc.OmitEmpty) && isEmpty(encoded) {
			continue
		}

		mapValue := encoded.Interface()

		// If ft should be squashed and v is a Data, then merge v into d.
		if ft.Squash {
			if data, ok := mapValue.(RawData); ok {
				for k, v := range data {
					d[k] = v
				}

				continue
			}
		}

		d[k] = mapValue
	}
}

func isEmpty(val reflect.Value) bool {
	if !val.IsValid() || val.IsZero() {
		return true
	}

	switch val.Kind() {
	case reflect.Map, reflect.Slice:
		return val.Len() == 0
	}

	return false
}

var valueOfEmptyString = reflect.ValueOf("")
var valueOfInvalid = reflect.ValueOf(nil)

func (enc *Encoder) encodeMapValue(val reflect.Value) reflect.Value {
	if !val.IsValid() {
		return valueOfInvalid
	}

	switch val.Type() {
	case typeOfTime:
		return val

	case typeOfDuration:
		if val.Int() == 0 {
			return valueOfEmptyString
		}

		method := val.MethodByName("String")
		rets := method.Call(nil)
		return rets[0]
	}

	switch val.Kind() {
	case reflect.Invalid:
		return valueOfInvalid

	case reflect.String:
		// Support json.Number. Convert it to a number if possible.
		if num, ok := val.Interface().(json.Number); ok {
			i64, err := num.Int64()

			if err == nil {
				return reflect.ValueOf(i64)
			}

			f64, err := num.Float64()

			if err == nil {
				return reflect.ValueOf(f64)
			}

			// Fallback to string.
			return reflect.ValueOf(string(num))
		}

		return val

	case reflect.Slice:
		return enc.copySlice(val)

	case reflect.Array:
		slice := val.Slice3(0, val.Len(), val.Cap())
		return enc.copySlice(slice)

	case reflect.Interface, reflect.Pointer:
		val = val.Elem()
		return enc.encodeMapValue(val)

	case reflect.Map:
		t := val.Type()
		kt := t.Key()

		// Unable to process map value other than map[string]any.
		if k := kt.Kind(); k != reflect.String {
			return clone.FromHeap().Clone(val)
		}

		d := RawData{}
		dval := reflect.ValueOf(d)

		if val.Len() == 0 {
			return dval
		}

		iter := val.MapRange()

		for iter.Next() {
			k := iter.Key()
			v := iter.Value()

			d[k.String()] = enc.encodeMapValue(v).Interface()
		}

		return dval

	case reflect.Struct:
		d := enc.encodeStruct(val)
		return reflect.ValueOf(d)

	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		// Always returns nil for these types.
		return valueOfInvalid
	}

	return val
}

func (enc *Encoder) copySlice(slice reflect.Value) reflect.Value {
	et := slice.Type().Elem()

	switch et.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		// No change for primitive types.
		return slice

	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		// Always returns a slice of nils for these types.
		return reflect.MakeSlice(slice.Type(), slice.Len(), slice.Cap())

	default:
		l := slice.Len()
		t := slice.Type()

		// Special case for []Data. Convert it to []RawData.
		if et.AssignableTo(typeOfData) {
			t = typeOfRawDataSlice
		}

		copied := reflect.MakeSlice(t, l, slice.Cap())

		for i := 0; i < l; i++ {
			val := enc.encodeMapValue(slice.Index(i))
			copied.Index(i).Set(val)
		}

		return copied
	}
}
