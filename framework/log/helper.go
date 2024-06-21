package log

import (
	"context"
	goLog "log"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
)

type Helper struct {
	*log.Helper
}

// NewHelper new a logger helper.
func NewHelper(logger log.Logger, opts ...log.Option) *Helper {
	lh := log.NewHelper(logger, opts...)
	return &Helper{Helper: lh}
}

func (h *Helper) GetKratosLogHelper() *log.Helper {
	return h.Helper
}

func (h *Helper) structMsg(ctx context.Context, level log.Level, keyvals ...interface{}) {
	h2 := h.WithContext(ctx)
	valMap := KeyvalsToMap(keyvals)
	h2.Log(level, "struct_msg", valMap)
}

func KeyvalsToMap(keyvals []interface{}) map[string]interface{} {
	if len(keyvals)%2 != 0 {
		goLog.Println("keyvals个数须为偶数")
	}

	valMap := map[string]interface{}{}
	for i := 1; i < len(keyvals); i += 2 {
		switch typeVal := keyvals[i-1].(type) {
		case string:
			valMap[typeVal] = keyvals[i]
		case int:
			pairKey := strconv.Itoa(typeVal)
			valMap[pairKey] = keyvals[i]
		}
	}
	return valMap
}

func (h *Helper) Debugc(ctx context.Context, a ...interface{}) {
	h2 := h.WithContext(ctx)
	h2.Debug(a...)
}

func (h *Helper) Infoc(ctx context.Context, a ...interface{}) {
	h2 := h.WithContext(ctx)
	h2.Info(a...)
}

func (h *Helper) Warnc(ctx context.Context, a ...interface{}) {
	h2 := h.WithContext(ctx)
	h2.Warn(a...)
}

func (h *Helper) Errorc(ctx context.Context, a ...interface{}) {
	h2 := h.WithContext(ctx)
	h2.Error(a...)
}

func (h *Helper) Fatalc(ctx context.Context, a ...interface{}) {
	h2 := h.WithContext(ctx)
	h2.Fatal(a...)
}

func (h *Helper) Debugs(ctx context.Context, keyvals ...interface{}) {
	h.structMsg(ctx, log.LevelDebug, keyvals...)
}

func (h *Helper) Infos(ctx context.Context, keyvals ...interface{}) {
	h.structMsg(ctx, log.LevelInfo, keyvals...)
}

func (h *Helper) Warns(ctx context.Context, keyvals ...interface{}) {
	h.structMsg(ctx, log.LevelWarn, keyvals...)
}

func (h *Helper) Errors(ctx context.Context, keyvals ...interface{}) {
	h.structMsg(ctx, log.LevelError, keyvals...)
}

func (h *Helper) Fatals(ctx context.Context, keyvals ...interface{}) {
	h.structMsg(ctx, log.LevelFatal, keyvals...)
}

func (h *Helper) Debuga(ctx context.Context, keyvals ...interface{}) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.AppendLog(log.LevelDebug, keyvals...)
	}
}

func (h *Helper) Infoa(ctx context.Context, keyvals ...interface{}) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.AppendLog(log.LevelInfo, keyvals...)
	}
}

func (h *Helper) Warna(ctx context.Context, keyvals ...interface{}) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.AppendLog(log.LevelWarn, keyvals...)
	}
}

func (h *Helper) Errora(ctx context.Context, keyvals ...interface{}) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.AppendLog(log.LevelError, keyvals...)
	}
}

func (h *Helper) Fatala(ctx context.Context, keyvals ...interface{}) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.AppendLog(log.LevelFatal, keyvals...)
	}
}

func (h *Helper) Debugm(ctx context.Context, a ...interface{}) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.AppendFmtLog(log.LevelDebug, a...)
	}
}

func (h *Helper) Infom(ctx context.Context, a ...interface{}) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.AppendFmtLog(log.LevelInfo, a...)
	}
}

func (h *Helper) Warnm(ctx context.Context, a ...interface{}) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.AppendFmtLog(log.LevelWarn, a...)
	}
}

func (h *Helper) Errorm(ctx context.Context, a ...interface{}) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.AppendFmtLog(log.LevelError, a...)
	}
}

func (h *Helper) Fatalm(ctx context.Context, a ...interface{}) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.AppendFmtLog(log.LevelFatal, a...)
	}
}

func WithHelper(ctx context.Context, l log.Logger, kv ...interface{}) (lh *Helper) {
	lg := log.With(l, kv...)
	lh = NewHelper(lg)
	return
}

func WithModuleHelper(ctx context.Context, l log.Logger, module string) (lh *Helper) {
	lg := log.With(l, "module", module)
	lh = NewHelper(lg)
	return
}
