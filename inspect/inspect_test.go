package inspect

import "testing"

func emptyInOut() {}

func fiveInTenOut(a, b, c, d, e int) (int, int, int, int, int, int, int, int, int, int) {
	return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
}
func TestNumInArgs(t *testing.T) {
	var inExpected, inResult int

	inExpected = 0
	inResult = NumInArgs(emptyInOut)
	if inResult != inExpected {
		t.Fatalf("Arg count should be %d, has %d", inExpected, inResult)
	}

	inExpected = 5
	inResult = NumInArgs(fiveInTenOut)
	if inResult != inExpected {
		t.Fatalf("Arg count should be %d, has %d", inExpected, inResult)
	}
}

func TestNumOutArgs(t *testing.T) {
	var inExpected, inResult int

	inExpected = 0
	inResult = NumOutArgs(emptyInOut)
	if inResult != inExpected {
		t.Fatalf("Arg count should be %d, has %d", inExpected, inResult)
	}

	inExpected = 10
	inResult = NumOutArgs(fiveInTenOut)
	if inResult != inExpected {
		t.Fatalf("Arg count should be %d, has %d", inExpected, inResult)
	}
}

func TestNamedCall(t *testing.T) {
	var result, callResult bool

	result = false
	callResult = NamedCall(fiveInTenOut)
	if callResult != result {
		t.Fatalf("Arg count should be %v, has %v", callResult, result)
	}
}
