package log

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

//"request.operation", log.OperationValuer(),
//"http.method", log.HttpMethodValuer(),
//"http.url.path", log.UrlPathValuer(),
//"http.url.host", log.UrlHostValuer(),
//"duration", log.StatDurationValuer(),
type HttpValue struct {
	Method string `json:"method"`
	URL    *url.URL
	Host   string
}

func OperationValuer() log.Valuer {
	return func(ctx context.Context) interface{} {
		//transport断言
		tr, ok := transport.FromServerContext(ctx)
		if !ok {
			return ""
		}

		return tr.Operation()
	}
}

func ShortOperationValuer() log.Valuer {
	return func(ctx context.Context) interface{} {
		//transport断言
		tr, ok := transport.FromServerContext(ctx)
		if !ok {
			return ""
		}

		op := tr.Operation()
		lastIdx := strings.LastIndex(op, ".")
		return op[lastIdx+1:]
	}
}

func HttpMethodValuer() log.Valuer {
	return func(ctx context.Context) interface{} {
		//transport断言
		tr, ok := transport.FromServerContext(ctx)
		if !ok {
			return ""
		}

		var httpTr *khttp.Transport
		var httpReq *http.Request

		if tr.Kind() == transport.KindHTTP {
			httpTr = tr.(*khttp.Transport)
			httpReq = httpTr.Request()

			return httpReq.Method
		}
		return ""
	}
}

func UrlPathValuer() log.Valuer {
	return func(ctx context.Context) interface{} {
		//transport断言
		tr, ok := transport.FromServerContext(ctx)
		if !ok {
			return ""
		}

		var httpTr *khttp.Transport
		var httpReq *http.Request

		if tr.Kind() == transport.KindHTTP {
			httpTr = tr.(*khttp.Transport)
			httpReq = httpTr.Request()

			if httpReq.URL != nil {
				return httpReq.URL.Path
			}
		}
		return ""
	}
}

func UrlHostValuer() log.Valuer {
	return func(ctx context.Context) interface{} {
		//transport断言
		tr, ok := transport.FromServerContext(ctx)
		if !ok {
			return ""
		}

		var httpTr *khttp.Transport
		var httpReq *http.Request

		if tr.Kind() == transport.KindHTTP {
			httpTr = tr.(*khttp.Transport)
			httpReq = httpTr.Request()

			if httpReq.URL != nil {
				return httpReq.URL.Host
			}
		}
		return ""
	}
}

func HttpHostValuer() log.Valuer {
	return func(ctx context.Context) interface{} {
		//transport断言
		tr, ok := transport.FromServerContext(ctx)
		if !ok {
			return ""
		}

		var httpTr *khttp.Transport
		var httpReq *http.Request

		if tr.Kind() == transport.KindHTTP {
			httpTr = tr.(*khttp.Transport)
			httpReq = httpTr.Request()

			return httpReq.Host
		}
		return ""
	}
}
