package data

import (
	"reflect"
	"time"
)

var (
	typeOfObject       = reflect.TypeOf(RawData{})
	typeOfData         = reflect.TypeOf(Data{})
	typeOfTime         = reflect.TypeOf(time.Time{})
	typeOfDuration     = reflect.TypeOf(time.Duration(0))
	typeOfRawDataSlice = reflect.TypeOf([]RawData{})
)
