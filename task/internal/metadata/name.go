package metadata

import "context"

type contextKey string

var nameKey contextKey = "name"

func WithName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, nameKey, name)
}

func GetName(ctx context.Context) string {
	if v := ctx.Value(nameKey); v != nil {
		return v.(string)
	}
	return "<Anonymous Task>"
}
