package log

import (
	"context"
	"fmt"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
)

type ctxLogStackItem struct {
	//LevelNum          log.Level
	Level             string
	StepDuration      string
	TotalDuration     string
	stepTimeDuration  time.Duration
	totalTimeDuration time.Duration
	Keyvals           map[string]interface{}
}

type StackLog struct {
	startTime    time.Time
	ctxLogStack  []*ctxLogStackItem
	logLevel     klog.Level
	customModule string
}

func NewStackLog(logLevel klog.Level) *StackLog {
	return &StackLog{startTime: time.Now(), logLevel: logLevel}
}

func (sl *StackLog) AppendLog(level klog.Level, keyvals ...interface{}) {
	sl.appendLog(level, KeyvalsToMap(keyvals))
}

func (sl *StackLog) AppendFmtLog(level klog.Level, a ...interface{}) {
	sl.appendLog(level, map[string]interface{}{"msg": fmt.Sprint(a...)})
}

func (sl *StackLog) appendLog(level klog.Level, kvs map[string]interface{}) {
	if level < sl.logLevel {
		return
	}
	prevTotalDuration := time.Duration(0)
	if len(sl.ctxLogStack) > 0 {
		prevTotalDuration = sl.ctxLogStack[len(sl.ctxLogStack)-1].totalTimeDuration
	}
	notTime := time.Now()

	totalDuration := notTime.Sub(sl.startTime)
	stepDuration := totalDuration - prevTotalDuration

	sl.ctxLogStack = append(sl.ctxLogStack, &ctxLogStackItem{
		//LevelNum:          level,
		Level:             level.String(),
		Keyvals:           kvs,
		stepTimeDuration:  stepDuration,
		totalTimeDuration: totalDuration,
		StepDuration:      stepDuration.String(),
		TotalDuration:     totalDuration.String(),
	})
}

func (sl *StackLog) GetLogs() []*ctxLogStackItem {
	return sl.ctxLogStack
}

func (sl *StackLog) SetCustomModule(module string) {
	sl.customModule = module
}

func (sl *StackLog) GetCustomModule() string {
	return sl.customModule
}

type ctxStackLogKey struct{}

// NewServerContext creates a new context with client md attached.
func NewStackLogServerContext(ctx context.Context, sl *StackLog) context.Context {
	return context.WithValue(ctx, ctxStackLogKey{}, sl)
}

// FromServerContext returns the server metadata in ctx if it exists.
func CtxStackLogFromServerContext(ctx context.Context) (sl *StackLog, ok bool) {
	sl, ok = ctx.Value(ctxStackLogKey{}).(*StackLog)
	return
}

func WithStackLogContext(ctx context.Context, logLevel klog.Level) context.Context {
	stackLog := NewStackLog(logLevel)
	stackLog.AppendLog(klog.LevelInfo, "step", "begin")
	return NewStackLogServerContext(ctx, stackLog)
}

func InfoStackLogs(ctx context.Context, logger *klog.Helper) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.AppendLog(klog.LevelInfo, "step", "end")
		logger.Log(klog.LevelInfo, "stack_log", sl.GetLogs())
	}
}

func ChangeStackLogModule(ctx context.Context, module string) {
	sl, ok := CtxStackLogFromServerContext(ctx)
	if ok {
		sl.SetCustomModule(module)
	}
}
