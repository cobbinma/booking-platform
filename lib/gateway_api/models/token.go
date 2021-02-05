package models

import (
	"context"
	"fmt"
)

type ctxKey string

const tokenCtxKey ctxKey = "token-ctx-key"

func GetTokenFromCtx(ctx context.Context) (string, error) {
	if token, ok := ctx.Value(tokenCtxKey).(string); ok {
		return token, nil
	}

	return "", fmt.Errorf("could not get token from context")
}

func AddTokenToCtx(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenCtxKey, token)
}
