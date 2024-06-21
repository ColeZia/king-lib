package middlewares

import (
	"context"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"gl.king.im/king-lib/framework/internal/di"
	"gl.king.im/king-lib/framework/log"
)

func ContextLog() middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			logger := di.GetLogger()

			ctxLogger := klog.WithContext(ctx, logger)

			ctx = log.NewLoggerServerContext(ctx, ctxLogger)

			rsp, err := next(ctx, req)

			return rsp, err
		}
	}
}
