package filters

import (
	netHttp "net/http"

	"gl.king.im/king-lib/framework/internal/tracing"
)

var UseDebugInfoFilter bool

func DebugInfo(next netHttp.Handler) netHttp.Handler {
	UseDebugInfoFilter = true
	return netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
		ctx := tracing.NewDebugInfoContext(r.Context(), &[]*tracing.DebugInfo{})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
