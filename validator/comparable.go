package validator

// Comparable represents a data which can be compared with another data.
type Comparable[T any] interface {
	// Compare compares v with the data.
	//
	//     - Returns > 0 if the data is greater than v;
	//     - Returns < 0 if it is less than v;
	//     - Returns 0 if it is equal to v.
	Compare(v T) int
}
