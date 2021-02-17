package object

import "testing"

func TestAtomicSetValue(t *testing.T) {
	a := &Atomic{}
	if err := a.SetValue(1, UntypedAtomicType); err != nil {
		t.Errorf(err.Error())
	}
	if a.Type() != UntypedAtomicType {
		t.Errorf("got=%s, expected=%s", a.Type(), UntypedAtomicType)
	}
	if err := a.SetValue(128, ByteType); err == nil {
		t.Errorf("byte type cannot have %d value", 128)
	}
	if err := a.SetValue(128, ShortType); err != nil {
		t.Errorf("got=%s, expected=%s", a.Type(), ShortType)
	}
}

// func TestQName(t *testing.T) {
// 	q := &QName{}
// 	q.Local.SetValue()
// }
