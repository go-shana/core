package data

import (
	"reflect"

	"github.com/huandu/go-clone"
)

// Merge merges multiple data from left to right to d,
// and if there are same keys, they will be merged deeply.
//
// Merging strategy:
//
//   - For Data type value, merge same key deeply.
//   - If same key exists and value types are the same, both are Data or slice, merge value deeply.
//   - If same key exists and value types are different, later value overwrites previous value.
//   - For slice type value, if two slices are the same type, later slice values will be appended.
func Merge(data ...Data) (d Data) {
	if len(data) == 0 {
		return emptyData
	}

	target := RawData{}
	merge(reflect.ValueOf(target), data[0].data, data[1:]...)
	return Data{
		data: target,
	}
}

// MergeTo merges multiple data from left to right to the target,
// and if there are same keys, they will be merged deeply.
// If target is nil, it will be ignored and no data will be merged.
//
// More details about merging strategy, please refer to `Merge`'s document.
func MergeTo(target *Data, data ...Data) {
	if target == nil || len(data) == 0 {
		return
	}

	if target.data == nil {
		target.data = RawData{}
	}

	merge(reflect.ValueOf(target.data), data[0].data, data[1:]...)
}

func merge(target reflect.Value, data RawData, remaining ...Data) {
	for k, v := range data {
		key := reflect.ValueOf(k)
		from := target.MapIndex(key)
		to := mergeValue(from, v)

		target.SetMapIndex(key, to)
	}

	if len(remaining) == 0 {
		return
	}

	merge(target, remaining[0].data, remaining[1:]...)
}

// mergeValue assumes that target and v are values in Data, so there will be no special types
// like ptr, struct, interface, and all map types are Data.
func mergeValue(target reflect.Value, v any) reflect.Value {
	if v == nil {
		return target
	}

	data := reflect.ValueOf(v)

	if target.IsValid() {
		for target.Kind() == reflect.Interface {
			target = target.Elem()
		}

		if target.Type() == data.Type() {
			switch target.Kind() {
			case reflect.Map:
				if target.IsNil() {
					target = reflect.MakeMap(target.Type())
				}

				iter := data.MapRange()

				for iter.Next() {
					key := iter.Key()
					from := target.MapIndex(key)
					to := mergeValue(from, iter.Value().Interface())

					target.SetMapIndex(key, to)
				}

				return target

			case reflect.Slice:
				return reflect.AppendSlice(target, data)
			}
		}
	}

	return reflect.ValueOf(clone.Clone(v))
}
