package tracing

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

type debugInfoCtxKey struct{}

// NewServerContext creates a new context with client md attached.
func NewDebugInfoContext(ctx context.Context, debugInfo interface{}) context.Context {
	return context.WithValue(ctx, debugInfoCtxKey{}, debugInfo)
}

// FromServerContext returns the server metadata in ctx if it exists.
func DebugInfoFromServerContext(ctx context.Context) *[]*DebugInfo {
	di := ctx.Value(debugInfoCtxKey{})
	d, ok := di.(*[]*DebugInfo)
	if !ok {
		log.Println("DebugInfoFromServerContext:DebugInfo asserting not ok")
	}

	return d
}

type DebugInfo struct {
	Time          time.Time
	Duration      time.Duration //阶段耗时-单位为毫秒--未考虑并发情况
	TotalDuration time.Duration //总体耗时-单位为毫秒--未考虑并发情况
	Data          interface{}
}

//func (d *DebugInfo) AppendTraceDebugInfo(data interface{}) {
//	d.Data = data
//}

func AppendTraceDebugInfo(ctx context.Context, data interface{}) {
	//由于kratos框架的encoder无法直接拿到业务代码或中间件的context，只能通过http Request的context来间接进行context的数据传递
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		return
	}

	trKind := tr.Kind()
	var httpTr *khttp.Transport
	var httpReq *http.Request

	if trKind == transport.KindHTTP {
		httpTr = tr.(*khttp.Transport)
		httpReq = httpTr.Request()
		d := DebugInfoFromServerContext(httpReq.Context())
		if d == nil {
			return
		}
		dv := *d
		var prevItem *DebugInfo
		if len(dv) > 0 {
			prevItem = dv[len(dv)-1]
		}

		timeNow := time.Now()
		newItem := &DebugInfo{
			Time: timeNow,
			Data: data,
		}

		if prevItem != nil {
			newItem.Duration = timeNow.Sub(prevItem.Time) / time.Millisecond
		}

		*d = append(*d, newItem)
		//d.AppendTraceDebugInfo(data)

	}

	//d := DebugInfoFromServerContext(ctx)
	//d.AppendTraceDebugInfo(data)
}
