package jsonrpc2

import "testing"

func Test_callsAdd(t *testing.T) {
	c := make(calls)
	c.Add("name1", 1)
	c.Add("name2", 1)

	if len(c) != 2 {
		t.Fatalf("Calls Error code is incorrect %d != 2", len(c))
	}
}

func Test_callsAddGet(t *testing.T) {
	c := make(calls)
	c.Add("name1", 1)

	value, _ := c.Get("name1")
	result := value.(int)
	// iAreaId, ok := value.(int)

	if result != 1 {
		t.Fatalf("Calls should return 1, not %d", result)
	}
}
