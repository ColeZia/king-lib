package tracing

import "context"

type globalTraceInfoCtxKey struct{}

type GlobalTraceInfo struct {
	IsUserBeginNode bool
}

// NewServerContext creates a new context with client md attached.
func NewGlobalTraceInfoServerContext(ctx context.Context, info *GlobalTraceInfo) context.Context {
	return context.WithValue(ctx, globalTraceInfoCtxKey{}, info)
}

// FromServerContext returns the server metadata in ctx if it exists.
func GlobalTraceInfoFromServerContext(ctx context.Context) (info *GlobalTraceInfo, ok bool) {
	info, ok = ctx.Value(globalTraceInfoCtxKey{}).(*GlobalTraceInfo)
	return
}
