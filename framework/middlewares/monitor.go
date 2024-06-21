package middlewares

import (
	"context"
	"strconv"
	"time"

	skyProm "gl.king.im/goserver/sky-agent/v2/lib/client_golang/prometheus"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func MonitorMiddleware(requests *skyProm.CounterVec, seconds *skyProm.HistogramVec) func(handler middleware.Handler) middleware.Handler {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			if requests != nil {
				requests.With(skyProm.Labels{
					"kind":      kind,
					"operation": operation,
					"code":      strconv.Itoa(code),
					"reason":    reason,
				}).Inc()
			}
			if seconds != nil {
				seconds.With(skyProm.Labels{
					"kind":      kind,
					"operation": operation,
				}).Observe(time.Since(startTime).Seconds())
			}
			return reply, err
		}
	}
}
