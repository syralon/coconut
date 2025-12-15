package errors

import "testing"

func f1() {
	panic("hello")
}

func TestRecovery(t *testing.T) {
	defer func() {
		err := Recovery(recover())
		t.Log(err)
		t.Log(err.stack)
	}()
	f1()
}
