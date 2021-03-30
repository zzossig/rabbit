package util

import "testing"

func TestNCName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"xs:string", false},
		{"xs:string()", false},
		{"string", true},
		{"string!", false},
		{"Q{http://example.com/ns}invoice", false},
	}

	for _, tt := range tests {
		b := IsNCName(tt.input)
		if b != tt.expected {
			t.Errorf("got=%v, expected=%v", b, tt.expected)
		}
	}
}

func TestQName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"xs:string", true},
		{"xs:string()", false},
		{"string", true},
		{"string!", false},
		{"Q{http://example.com/ns}invoice", false},
	}

	for _, tt := range tests {
		b := IsQName(tt.input)
		if b != tt.expected {
			t.Errorf("got=%v, expected=%v", b, tt.expected)
		}
	}
}

func TestEQName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"xs:string", true},
		{"xs:string()", false},
		{"string", true},
		{"string!", false},
		{"Q{http://example.com/ns}invoice", true},
		{"{http://example.com/ns}invoice", false},
		{"Q{http://example.com/ns}", false},
	}

	for _, tt := range tests {
		b := IsEQName(tt.input)
		if b != tt.expected {
			t.Errorf("got=%v, expected=%v", b, tt.expected)
		}
	}
}

func TestIsReverseAxis(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"abcd::", false},
		{"child::", false},
		{"descendant::", false},
		{"descendant-or-self::", false},
		{"attribute::", false},
		{"self::", false},
		{"following-sibling::", false},
		{"following::", false},
		{"namespace::", false},
		{"ppparent::", false},
		{"parentparent::", false},
		{"parent::", true},
		{"ancestor::", true},
		{"preceding-sibling::", true},
		{"preceding::", true},
		{"ancestor-or-self::", true},
	}

	for _, tt := range tests {
		b := IsReverseAxis(tt.input)
		if b != tt.expected {
			t.Errorf("got=%v, expected=%v", b, tt.expected)
		}
	}
}

func TestIsForwardAxis(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"abcd::", false},
		{"abchild::", false},
		{"childchild::", false},
		{"child::", true},
		{"descendant::", true},
		{"descendant-or-self::", true},
		{"attribute::", true},
		{"self::", true},
		{"following-sibling::", true},
		{"following::", true},
		{"namespace::", true},
		{"parent::", false},
		{"ancestor::", false},
		{"preceding-sibling::", false},
		{"preceding::", false},
		{"ancestor-or-self::", false},
	}

	for _, tt := range tests {
		b := IsForwardAxis(tt.input)
		if b != tt.expected {
			t.Errorf("got=%v, expected=%v", b, tt.expected)
		}
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"612", true},
		{"012", true},
		{"12.1", false},
		{"12.1e-1", false},
		{"12.1e+1", false},
	}

	for _, tt := range tests {
		b := IsDigit(tt.input)
		if b != tt.expected {
			t.Errorf("got=%v, expected=%v", b, tt.expected)
		}
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{".05e2", true},
		{".05e-2", true},
		{".05e+2", true},
		{".05e?2", false},
		{".05e.2", false},
		{".05E2", true},
		{".05E-2", true},
		{".05E+2", true},
		{".05E?2", false},
		{".05E.2", false},
		{"12.05e1", true},
		{"12.05e11", true},
		{"12.05e11-1", false},
		{"12.05e11.1", false},
		{"12.01", true},
	}

	for _, tt := range tests {
		b := IsNumber(tt.input)
		if b != tt.expected {
			t.Errorf("got=%v, expected=%v", b, tt.expected)
		}
	}
}

func TestIsNodeComp(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"<<", true},
		{"<<<", false},
		{">>", true},
		{">>>", false},
		{"is", true},
		{"isis", false},
		{"iis", false},
		{"iss", false},
	}

	for _, tt := range tests {
		b := IsNodeComp(tt.input)
		if b != tt.expected {
			t.Errorf("got=%v, expected=%v", b, tt.expected)
		}
	}
}

func TestIsGeneralComp(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"=", true},
		{"==", false},
		{"!=", true},
		{"<=", true},
		{">=", true},
		{">", true},
		{"<", true},
		{"<<", false},
		{">>", false},
	}

	for _, tt := range tests {
		b := IsGeneralComp(tt.input)
		if b != tt.expected {
			t.Errorf("got=%v, expected=%v", b, tt.expected)
		}
	}
}

func TestIsValueComp(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"eq", true},
		{"ne", true},
		{"ab", false},
		{"cd", false},
		{"lt", true},
		{"le", true},
		{"gt", true},
		{"ge", true},
		{"gee", false},
		{">>", false},
	}

	for _, tt := range tests {
		b := IsValueComp(tt.input)
		if b != tt.expected {
			t.Errorf("got=%v, expected=%v", b, tt.expected)
		}
	}
}
