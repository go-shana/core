package data

import (
	"fmt"
	"testing"

	"github.com/huandu/go-assert"
)

func ExamplePatch() {
	patch := NewPatch()

	// Delete d["v2"]、d["v3"][1]、d["v4"]["v4-4"]。
	patch.Add([]string{"v2", "v3.1", "v4.v4-1"}, nil)

	// Add data.
	patch.Add(nil, map[string]Data{
		// Add data at root.
		"": Make(RawData{
			"v1": []int{2, 3},
			"v2": 456,
		}),

		// Add data in d["v4"].
		"v4": Make(RawData{
			"v4-1": "new",
		}),
	})

	// Delete and add data at the same time.
	patch.Add([]string{"v4.v4-2"}, map[string]Data{
		"v4": Make(RawData{
			"v4-2": RawData{
				"new": true,
			},
		}),
	})

	d := Make(RawData{
		"v1": []int{1},
		"v2": 123,
		"v3": []string{"first", "second", "third"},
		"v4": RawData{
			"v4-1": "old",
			"v4-2": RawData{
				"old": true,
			},
		},
	})
	patch.ApplyTo(&d)
	fmt.Println(d)

	// Output:
	// {"v1":[1,2,3],"v2":456,"v3":["first","third"],"v4":{"v4-1":"new","v4-2":{"new":true}}}
}

func TestPatch(t *testing.T) {
	cases := []struct {
		Deletes  []string
		Updates  map[string]Data
		Target   Data
		Result   Data
		HasError bool
	}{
		{ // Do nothing.
			nil,
			nil,
			fullData,
			fullData,
			false,
		},
		{ // Delete all data.
			[]string{""},
			nil,
			fullData,
			Data{},
			false,
		},
		{ // Delete array element.
			[]string{"arr.1"},
			nil,
			Make(RawData{
				"arr": []int{1, 2, 3},
			}),
			Make(RawData{
				"arr": []int{1, 3},
			}),
			false,
		},
		{ // Add map element.
			nil,
			map[string]Data{
				"": Make(RawData{
					"v2": 222,
					"v3": RawData{
						"arr": []int{3, 4},
					},
				}),
				"v3.map": Make(RawData{
					"k1": false,
					"k2": 999,
				}),
				"v4.2.1.2": Make(RawData{
					"v4-1": "new",
					"v4-2": 2222,
				}),
			},
			Make(RawData{
				"v1": 1,
				"v2": 2,
				"v3": RawData{
					"arr": []int{1, 2},
					"map": RawData{
						"k1": true,
					},
				},
				"v4": [][][]RawData{
					nil, nil,
					{
						nil,
						{
							nil, nil,
							{
								"v4-1": "old",
							},
						},
					},
				},
			}),
			Make(RawData{
				"v1": 1,
				"v2": 222,
				"v3": RawData{
					"arr": []int{1, 2, 3, 4},
					"map": RawData{
						"k1": false,
						"k2": 999,
					},
				},
				"v4": [][][]RawData{
					nil, nil,
					{
						nil,
						{
							nil, nil,
							{
								"v4-1": "new",
								"v4-2": 2222,
							},
						},
					},
				},
			}),
			false,
		},
		{ // Delete and update at the same time.
			[]string{"sub_type", "anonymous_type.data_list.1"},
			map[string]Data{
				"": Make(RawData{
					"string": "xyz",
					"sub_type": RawData{
						"int32": int32(-32000),
					},
				}),
				"anonymous_type.data_list.1": Make(RawData{
					"a": "aaaa",
					"b": "bbbb",
				}),
			},
			fullData,
			func() Data {
				d := fullData.Clone()

				delete(d.data, "sub_type")
				slice := d.data["anonymous_type"].(RawData)["data_list"].([]RawData)
				slice[1] = slice[2]
				slice = slice[:2]
				d.data["anonymous_type"].(RawData)["data_list"] = slice

				d.data["string"] = "xyz"
				d.data["sub_type"] = RawData{
					"int32": int32(-32000),
				}
				dataList := d.data["anonymous_type"].(RawData)["data_list"].([]RawData)[1]
				dataList["a"] = "aaaa"
				dataList["b"] = "bbbb"

				return d
			}(),
			false,
		},
		{ // Delete with invalid path.
			[]string{"arr.3", "foo"},
			nil,
			Make(RawData{
				"arr": []int{1, 2, 3},
			}),
			Make(RawData{
				"arr": []int{1, 2, 3},
			}),
			false,
		},
		{ // Update invalid path.
			nil,
			map[string]Data{
				"foo": Make(RawData{
					"bar": 1,
				}),
			},
			Make(RawData{
				"arr": []int{1, 2, 3},
			}),
			Data{},
			true,
		},
		{ // Update data with invalid type.
			nil,
			map[string]Data{
				"foo": Make(RawData{
					"bar": 1,
				}),
			},
			Make(RawData{
				"foo": 123,
			}),
			Data{},
			true,
		},
	}
	a := assert.New(t)

	for i, c := range cases {
		target := c.Target.Clone()
		a.Use(&i, &c, &target)

		p := NewPatch()
		p.Add(c.Deletes, c.Updates)

		a.Equal(p.Actions(), []*PatchAction{
			{
				Deletes: c.Deletes,
				Updates: c.Updates,
			},
		})

		applied, err := p.Apply(c.Target)

		if c.HasError {
			a.NonNilError(err)
		} else {
			a.NilError(err)
			a.Equal(c.Result, applied)
		}

		err = p.ApplyTo(&target)

		if c.HasError {
			a.NonNilError(err)
		} else {
			a.NilError(err)
			a.Equal(c.Result, target)
		}

	}
}
