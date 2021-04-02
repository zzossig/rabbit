package rabbit

import (
	"testing"
)

func TestXPath(t *testing.T) {
	data := New().SetDoc("./eval/testdata/company_2.xml").Eval("//employee").Data()
	if len(data) != 5 {
		t.Errorf("result length should be 5. got=%d", len(data))
	}

	x := New()
	data2 := x.Eval("1+1").Data()
	if data2[0] != 2 {
		t.Errorf("result value should be 2. got=%d", data2[0])
	}
}
func BenchmarkXPath(b *testing.B) {
	x := New().SetDoc("./eval/testdata/company_2.xml")
	for n := 0; n < b.N; n++ {
		x.Eval("//employee").Data()
	}
}
