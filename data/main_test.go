package data

import "time"

type AllValue struct {
	Empty string `test:"empty,omitempty"`

	Bool           bool          `test:"bool"`
	Int            int           `test:"int"`
	Uint           uint          `test:"uint"`
	Time           time.Time     `test:"time"`
	Duration       time.Duration `test:"duration"`
	String         string        `test:"string"`
	FakeQuery      float64       `test:"fake.query"`
	Skipped        string        `test:"-"`
	SubType        SubType       `test:"sub_type"`
	*AnonymousType `test:"anonymous_type"`
	SquashType     *SquashType `test:"squash_type,squash"`
}

type SubType struct {
	Int8    int8     `test:"int8"`
	Int16   int16    `test:"int16"`
	Int32   int32    `test:"int32"`
	Int64   int64    `test:"int64"`
	Strings []string `test:"strings"`
}

type AnonymousType struct {
	Data       Data       `test:"data"`
	DataList   []Data     `test:"data_list"`
	Complex64  complex64  `test:"complex64"`
	Complex128 complex128 `test:"complex128"`
}

type SquashType struct {
	Uint8  uint8   `test:"uint8"`
	Uint16 uint16  `test:"uint16"`
	Uint32 uint32  `test:"uint32"`
	Uint64 uint64  `test:"uint64"`
	Uints  [6]uint `test:"uints"`
}

var (
	testTime = time.Date(2023, 2, 23, 14, 15, 16, 0, time.Local)

	fullData = Make(RawData{
		"bool":       true,
		"int":        -123,
		"uint":       uint(456),
		"time":       testTime,
		"duration":   "8.321s",
		"string":     "abcd",
		"fake.query": 67.89,
		"sub_type": RawData{
			"int8":    int8(-8),
			"int16":   int16(-16),
			"int32":   int32(-32),
			"int64":   int64(-64),
			"strings": []string{"this", "is", "", "string"},
		},
		"anonymous_type": RawData{
			"data": RawData{
				"foo": int64(123),
			},
			"data_list": []RawData{
				{"a": int64(1)},
				{},
				{"b": true},
			},
			"complex64":  complex64(complex(34, 5.5)),
			"complex128": complex(78.9, -10),
		},

		// Squashed
		"uint8":  uint8(8),
		"uint16": uint16(16),
		"uint32": uint32(32),
		"uint64": uint64(64),
		"uints":  []uint{4, 3, 2, 1, 0, 0},
	})
	allValues = &AllValue{
		Bool:      true,
		Int:       -123,
		Uint:      456,
		Time:      testTime,
		Duration:  8*time.Second + 321*time.Millisecond,
		String:    "abcd",
		FakeQuery: 67.89,
		SubType: SubType{
			Int8:    -8,
			Int16:   -16,
			Int32:   -32,
			Int64:   -64,
			Strings: []string{"this", "is", "", "string"},
		},
		AnonymousType: &AnonymousType{
			Data: Make(RawData{
				"foo": int64(123),
			}),
			DataList: []Data{
				Make(RawData{"a": int64(1)}),
				{},
				Make(RawData{"b": true}),
			},
			Complex64:  complex(34, 5.5),
			Complex128: complex(78.9, -10),
		},
		SquashType: &SquashType{
			Uint8:  8,
			Uint16: 16,
			Uint32: 32,
			Uint64: 64,
			Uints:  [6]uint{4, 3, 2, 1},
		},
	}
)
