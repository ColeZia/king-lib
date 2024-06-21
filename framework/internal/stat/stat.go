package stat

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type ctxStatKey struct{}

type CtxStat struct {
	StartTime time.Time
}

func NewCtxStat() *CtxStat {
	return &CtxStat{StartTime: time.Now()}
}

// NewServerContext creates a new context with client md attached.
func NewCtxStatServerContext(ctx context.Context, sl *CtxStat) context.Context {
	return context.WithValue(ctx, ctxStatKey{}, sl)
}

// FromServerContext returns the server metadata in ctx if it exists.
func CtxStatFromServerContext(ctx context.Context) (sl *CtxStat, ok bool) {
	sl, ok = ctx.Value(ctxStatKey{}).(*CtxStat)
	return
}

func StatDurationStrValuer() log.Valuer {
	return func(ctx context.Context) interface{} {
		ctxStatIns, ok := CtxStatFromServerContext(ctx)
		if !ok {
			return ""
		}

		return time.Now().Sub(ctxStatIns.StartTime).String()
	}
}
