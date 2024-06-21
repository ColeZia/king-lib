package tracing

import (
	"context"

	"gl.king.im/king-lib/framework/internal/tracing"
)

func AppendTraceDebugInfo(ctx context.Context, data interface{}) {
	tracing.AppendTraceDebugInfo(ctx, data)
}
