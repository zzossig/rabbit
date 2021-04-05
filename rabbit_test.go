package rabbit

import (
	"net/http"
	"testing"
)

func TestXPath(t *testing.T) {
	locations := New().SetDoc("./eval/testdata/company_2.xml").Eval("//company").Eval("./office").Evals("./@location")
	if len(locations) != 2 {
		t.Errorf("locations should have 2 xpath object")
	}
	if locations[0].Get() != "Seoul" || locations[1].Get() != "Busan" {
		t.Errorf("wrong attribute value")
	}

	items := New().SetDoc("./eval/testdata/company_2.xml").Evals("//company//age")
	mapping := []string{"25", "30", "30", "34", "44"}
	for i, item := range items {
		if item.Get() != mapping[i] {
			t.Errorf("expected=%s, got=%s", mapping[i], item.Get())
		}
	}

	items2 := New().SetDoc("./eval/testdata/company_2.xml").Eval("//company//age").GetAll()
	for i, item := range items2 {
		if item != mapping[i] {
			t.Errorf("expected=%s, got=%s", mapping[i], item)
		}
	}

	nodes := New().SetDoc("./eval/testdata/company_2.xml").Eval("//employee").Eval("./age").NodeAll()
	if len(nodes) != 5 {
		t.Errorf("result length should be 5. got=%d", len(nodes))
	}

	node2 := New().SetDocS("<h2>Hello, World!</h2>").Eval("//h2/text()").Node()
	if node2.Data != "Hello, World!" {
		t.Errorf("expected='Hello, World!', got=%s", node2.Data)
	}

	resp, err := http.Get("https://golang.org/")
	if err != nil {
		t.Errorf(err.Error())
	}
	node := New().SetDocR(resp).Eval("//title/text()").Node()
	if node.Data != "The Go Programming Language" {
		t.Errorf("unexpected node data. got=%s, expected='The Go Programming Language'", node.Data)
	}

	x := New()
	data2 := x.Eval("1+1").DataAll()
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
