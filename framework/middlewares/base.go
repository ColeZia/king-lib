package middlewares

import (
	"github.com/go-kratos/kratos/v2/middleware"
)

func BaseMiddleware(next middleware.Handler) middleware.Handler {
	return BaseMiddleWareCore(next, FrameworkMwCnf{
		Type: "kratos",
	})
}
