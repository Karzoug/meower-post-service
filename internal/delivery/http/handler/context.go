package handler

import (
	"context"

	"github.com/Karzoug/meower-post-service/internal/delivery/http/middleware"
)

func getUsernameFromContext(ctx context.Context) string {
	username, ok := ctx.Value(middleware.AuthUsernameKey).(string)
	if !ok {
		return ""
	}
	return username
}
