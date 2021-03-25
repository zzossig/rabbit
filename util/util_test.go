package util

import "testing"

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
