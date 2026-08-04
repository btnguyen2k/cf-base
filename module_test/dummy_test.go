package module_test

import "testing"

func TestDummy(t *testing.T) {
	testName := "TestDummy"
	t.Logf("%s started", testName)
	t.Logf("%s completed, passed", testName)
}
