package context

import (
	"context"
	"myphoto/models"
)

type contextKey string

const userKey contextKey = "user"

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	if value := ctx.Value(userKey); value != nil {
		if user, ok := value.(*models.User); ok {
			return user
		}
	}
	return nil
}
