package initer

import (
	"context"
	"testing"

	"github.com/go-shana/core/errors"
	"github.com/huandu/go-assert"
)

type testInitInt int

const expectedTestInitInt = 1

func (i *testInitInt) Init(ctx context.Context) {
	*i = expectedTestInitInt
}

var _ Initer = new(testInitInt)

type testInitPlainStruct struct {
	A int
	B string
	C *float64 // Should not be initialized automatically.
}

var _ Initer = new(testInitPlainStruct)

const (
	expectedTestInitPlainStructA = 2
	expectedTestInitPlainStructB = "expectedTestInitPlainStructB"
)

func (s *testInitPlainStruct) Init(ctx context.Context) {
	s.A = expectedTestInitPlainStructA
	s.B = expectedTestInitPlainStructB
}

type testInitNestedListStruct struct {
	Next *testInitNestedListStruct

	Disabled *testInitNestedStruct `init:"disabled"`
	private  *testInitNestedStruct
}

type testInitNestedStruct struct {
	A    string
	List testInitNestedListStruct
	Next *testInitNestedStruct `init:"disabled"`
}

var _ Initer = new(testInitNestedStruct)

const expectedTestInitNestedStructA = "expectedTestInitNestedStructA"

func (s *testInitNestedStruct) Init(ctx context.Context) {
	s.A = expectedTestInitNestedStructA
}

type testInitCannotSetStruct struct {
	A int
}

var _ Initer = new(testInitCannotSetStruct)

const expectedTestInitCannotSetStructA = 3

func (s testInitCannotSetStruct) Init(ctx context.Context) {
	s.A = expectedTestInitCannotSetStructA // Should not work right.
}

type testInitMap map[string]int

var _ Initer = testInitMap{}

const (
	expectedTestInitMapKey   = "foo"
	expectedTestInitMapValue = 123
)

func (s testInitMap) Init(ctx context.Context) {
	s[expectedTestInitMapKey] = expectedTestInitMapValue
}

type testInitComplexStruct struct {
	A string

	B1 testInitInt
	B2 ****testInitInt
	C1 testInitPlainStruct
	C2 *testInitPlainStruct `init:"omitempty"`
	D1 testInitNestedStruct
	D2 *testInitNestedStruct
	E1 testInitCannotSetStruct
	E2 *testInitCannotSetStruct
	F1 testInitMap

	private    testInitInt // Should not be initialized.
	NotTouched int         // Should not be touched by Init.
}

var _ Initer = new(testInitComplexStruct)

const (
	expectedTestInitComplexStructA          = "expectedTestInitComplexStructA"
	expectedTestInitComplexStructNotTouched = 4
)

func (s *testInitComplexStruct) Init(ctx context.Context) {
	s.A = expectedTestInitComplexStructA
}

type testInitFailedStruct struct {
	A string
}

var _ Initer = new(testInitFailedStruct)

var errTestInitFailed = errors.New("expected failure")

func (s *testInitFailedStruct) Init(ctx context.Context) {
	errors.Throw(errTestInitFailed)
}

func TestInit(t *testing.T) {
	// Defines several struct types which implements Initer.
	// Call Init to check whether all data is initialized as expected.
	// Covers all code branches as much as possible.

	a := assert.New(t)
	ctx := context.Background()

	var st *testInitComplexStruct
	a.NilError(Init(ctx, &st))

	a.Assert(st != nil)
	a.Equal(st.A, expectedTestInitComplexStructA)
	a.Assert(st.private == 0)

	a.Assert(st.B1 == expectedTestInitInt)
	a.Assert(st.B2 != nil)
	a.Assert(*st.B2 != nil)
	a.Assert(**st.B2 != nil)
	a.Assert(***st.B2 != nil)
	a.Assert(****st.B2 == expectedTestInitInt)

	a.Equal(st.C1.A, expectedTestInitPlainStructA)
	a.Equal(st.C1.B, expectedTestInitPlainStructB)
	a.Assert(st.C1.C == nil)
	a.Assert(st.C2 == nil)

	a.Equal(st.D1.A, expectedTestInitNestedStructA)
	a.Assert(st.D1.Next == nil)
	a.Assert(st.D1.List.Next == nil)
	a.Assert(st.D2 != nil)
	a.Equal(st.D2.A, expectedTestInitNestedStructA)
	a.Assert(st.D2.Next == nil)
	a.Assert(st.D2.List.Next == nil)

	a.Assert(st.F1 != nil)
	a.Equal(len(st.F1), 1)
	a.Equal(st.F1[expectedTestInitMapKey], expectedTestInitMapValue)

	a.NotEqual(st.E1.A, expectedTestInitCannotSetStructA)
	a.Assert(st.E2 != nil)
	a.NotEqual(st.E2.A, expectedTestInitCannotSetStructA)

	st.NotTouched = expectedTestInitComplexStructNotTouched
	a.NilError(Init(ctx, &st))
	a.Equal(st.NotTouched, expectedTestInitComplexStructNotTouched)
}

func TestInitWithError(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()

	var st *testInitFailedStruct
	a.NonNilError(Init(ctx, st))

	st = &testInitFailedStruct{}
	a.NonNilError(Init(ctx, st))
}
