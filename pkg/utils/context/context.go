package context

import "context"

type key struct{}

func SetPayload(ctx context.Context, payload interface{}) context.Context {
	return context.WithValue(ctx, key{}, payload)
}

func GetPayload(ctx context.Context) interface{} {
	return ctx.Value(key{})
}
