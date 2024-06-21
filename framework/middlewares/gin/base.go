package gin

import (
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"gl.king.im/king-lib/framework/middlewares"
)

var (
	ContextDeadlineDuration time.Duration = time.Millisecond * 1000 * 10 //单位为毫秒，默认为10秒
)

func Base(next middleware.Handler) middleware.Handler {
	return middlewares.BaseMiddleWareCore(next, middlewares.FrameworkMwCnf{
		Type:                    "gin",
		ContextDeadlineDuration: ContextDeadlineDuration,
	})
}
