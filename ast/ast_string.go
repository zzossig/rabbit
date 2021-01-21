package ast

import "github.com/zzossig/xpath/util"

// EQName ::= QName | URIQualifiedName
type EQName struct {
	value string
}

// Value is a getter for the value field
func (eqn *EQName) Value() string {
	return eqn.value
}

// SetValue is a setter for the value field
func (eqn *EQName) SetValue(name string) {
	if util.IsEQName(name) {
		eqn.value = name
	} else {
		// TODO occur error
	}
}

// NCName ::= Name - (Char* ':' Char*)
type NCName struct {
	value string
}

// Value is a getter for the value field
func (ncn *NCName) Value() string {
	return ncn.value
}

// SetValue is a setter for the value field
func (ncn *NCName) SetValue(name string) {
	if util.IsNCName(name) {
		ncn.value = name
	} else {
		// TODO occur error
	}
}
