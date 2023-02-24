package errors

// Check panics with the last error in values if it's not nil.
func Check(values ...any) {
	err := checkWithError(values)

	if err != nil {
		Throw(err)
	}
}

func checkWithError(values []any) error {
	l := len(values)

	if l == 0 {
		return nil
	}

	last := values[l-1]

	if last == nil {
		return nil
	}

	if cause, ok := last.(error); ok {
		return cause
	}

	return nil
}

// Check1 returns the in1 if err is nil and panics with err otherwise.
func Check1[T1 any](in1 T1, err error) (out1 T1) {
	if err != nil {
		Throw(err)
		return
	}

	out1 = in1
	return
}

// Check2 returns the in1 and in2 if err is nil and panics with err otherwise.
func Check2[T1, T2 any](in1 T1, in2 T2, err error) (out1 T1, out2 T2) {
	if err != nil {
		Throw(err)
		return
	}

	out1 = in1
	out2 = in2
	return
}

// Check3 returns the in1, in2 and in3 if err is nil and panics with err otherwise.
func Check3[T1, T2, T3 any](in1 T1, in2 T2, in3 T3, err error) (out1 T1, out2 T2, out3 T3) {
	if err != nil {
		Throw(err)
		return
	}

	out1 = in1
	out2 = in2
	out3 = in3
	return
}

// Check4 returns the in1, in2, in3 and in4 if err is nil and panics with err otherwise.
func Check4[T1, T2, T3, T4 any](in1 T1, in2 T2, in3 T3, in4 T4, err error) (out1 T1, out2 T2, out3 T3, out4 T4) {
	if err != nil {
		Throw(err)
		return
	}

	out1 = in1
	out2 = in2
	out3 = in3
	out4 = in4
	return
}

// Check5 returns the in1, in2, in3, in4 and in5 if err is nil and panics with err otherwise.
func Check5[T1, T2, T3, T4, T5 any](in1 T1, in2 T2, in3 T3, in4 T4, in5 T5, err error) (out1 T1, out2 T2, out3 T3, out4 T4, out5 T5) {
	if err != nil {
		Throw(err)
		return
	}

	out1 = in1
	out2 = in2
	out3 = in3
	out4 = in4
	out5 = in5
	return
}

// Check6 returns the in1, in2, in3, in4, in5 and in6 if err is nil and panics with err otherwise.
func Check6[T1, T2, T3, T4, T5, T6 any](in1 T1, in2 T2, in3 T3, in4 T4, in5 T5, in6 T6, err error) (out1 T1, out2 T2, out3 T3, out4 T4, out5 T5, out6 T6) {
	if err != nil {
		Throw(err)
		return
	}

	out1 = in1
	out2 = in2
	out3 = in3
	out4 = in4
	out5 = in5
	out6 = in6
	return
}

// Check7 returns the in1, in2, in3, in4, in5, in6 and in7 if err is nil and panics with err otherwise.
func Check7[T1, T2, T3, T4, T5, T6, T7 any](in1 T1, in2 T2, in3 T3, in4 T4, in5 T5, in6 T6, in7 T7, err error) (out1 T1, out2 T2, out3 T3, out4 T4, out5 T5, out6 T6, out7 T7) {
	if err != nil {
		Throw(err)
		return
	}

	out1 = in1
	out2 = in2
	out3 = in3
	out4 = in4
	out5 = in5
	out6 = in6
	out7 = in7
	return
}

// Check8 returns the in1, in2, in3, in4, in5, in6, in7 and in8 if err is nil and panics with err otherwise.
func Check8[T1, T2, T3, T4, T5, T6, T7, T8 any](in1 T1, in2 T2, in3 T3, in4 T4, in5 T5, in6 T6, in7 T7, in8 T8, err error) (out1 T1, out2 T2, out3 T3, out4 T4, out5 T5, out6 T6, out7 T7, out8 T8) {
	if err != nil {
		Throw(err)
		return
	}

	out1 = in1
	out2 = in2
	out3 = in3
	out4 = in4
	out5 = in5
	out6 = in6
	out7 = in7
	out8 = in8
	return
}
