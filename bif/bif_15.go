package bif

import "github.com/zzossig/rabbit/object"

func fnPosition(ctx *object.Context, args ...object.Item) object.Item {
	return NewInteger(ctx.CPos)
}

func fnLast(ctx *object.Context, args ...object.Item) object.Item {
	return NewInteger(ctx.CSize)
}
