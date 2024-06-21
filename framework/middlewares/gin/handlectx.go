package gin

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"

	kgin "github.com/go-kratos/gin"
)

func HandleGinCtx(next middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, replyErr error) {

		ginCtx, ok := kgin.FromGinContext(ctx)
		if ok {
			ginCtx.Request = ginCtx.Request.WithContext(ctx)
		}
		return next(ctx, req)
	}
}
