package coder

import (
	"encoding/json"
	"log"
	"net/http"
	netHttp "net/http"
	"strings"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"gl.king.im/king-lib/framework/alerting"
	"gl.king.im/king-lib/framework/alerting/feishu"
	"gl.king.im/king-lib/framework/alerting/types"
	eAlert "gl.king.im/king-lib/framework/errors/alert"
	"gl.king.im/king-lib/framework/internal/di"
	"gl.king.im/king-lib/framework/service/appinfo"

	ke "github.com/go-kratos/kratos/v2/errors"
)

type alertMsgTypeRequstURL struct {
	//Scheme      string
	//Opaque      string    // encoded opaque data
	//User        *Userinfo // username and password information
	Host string // host or host:port
	Path string // path (relative paths may omit leading slash)
	//RawPath string // encoded path hint (see EscapedPath method)
	//ForceQuery bool   // append a query ('?') even if RawQuery is empty
	RawQuery string // encoded query values, without '?'
	//Fragment    string // fragment for references, without '#'
	//RawFragment string // encoded fragment hint (see EscapedFragment method)
}

type alertMsgTypeRequst struct {
	Method string

	//URL *url.URL
	URL alertMsgTypeRequstURL

	//Proto string // "HTTP/1.0"
	//ProtoMajor int    // 1
	//ProtoMinor int    // 0

	Header http.Header

	//ContentLength int64
	//TransferEncoding []string

	Host string

	//Form     url.Values
	//PostForm url.Values
	//MultipartForm *multipart.Form

	//RemoteAddr string

	//RequestURI string
}

type alertMsgTypeInfo struct {
	ServiceName string
	TraceId     string
	Name        string
	Version     string
	Id          string
	//Env string
}
type alertMsgType struct {
	Info alertMsgTypeInfo

	HttpRequest alertMsgTypeRequst
}

//type BossError struct {
//	ke.Error
//	ErrCode string
//}

func ErrorEncoder() khttp.EncodeErrorFunc {
	return func(w netHttp.ResponseWriter, r *netHttp.Request, err error) {
		var (
			doAlert        bool
			alertMsg       string
			contentBuilder func(error) string
			alertingIns    *alerting.Alerting
			alertFun       alerting.AlertFun
		)

		alertingIns = di.Container.DefaultAlerting

		switch asVal := err.(type) {
		case *eAlert.AlertError:
			_ = asVal
			doAlert = true
			if asVal.Alerting != nil {
				alertingIns = asVal.Alerting
			}

			if asVal.AlertFun != nil {
				alertFun = asVal.AlertFun
			}

		default:

		}

		// 拿到error并转换成kratos Error实体
		se := ke.FromError(err)
		// 通过Request Header的Accept中提取出对应的编码器
		codec, _ := CodecForRequest(r, "Accept")

		body, err := codec.Marshal(se)
		if err != nil {
			w.WriteHeader(netHttp.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", ContentType(codec.Name()))
		// 设置HTTP Status Code
		errCode := int(se.Code)
		if errCode < 1000 {
			w.WriteHeader(errCode)
		}

		if errCode >= 500 && errCode <= 599 {
			doAlert = true
		}

		if doAlert {
			traceId := w.Header().Get("X-TRACE-ID")

			headerMsg := netHttp.Header{}
			for hk, hv := range r.Header {
				if http.CanonicalHeaderKey(hk) == "X-Md-Global-Call-Scenario" ||
					http.CanonicalHeaderKey(hk) == "X-Request-Id" ||
					http.CanonicalHeaderKey(hk) == "X-B3-Traceid" ||
					http.CanonicalHeaderKey(hk) == "Origin" ||
					http.CanonicalHeaderKey(hk) == "Referer" ||
					http.CanonicalHeaderKey(hk) == "Remoteip" ||
					http.CanonicalHeaderKey(hk) == "Referer" ||
					http.CanonicalHeaderKey(hk) == "User-Agent" {
					headerMsg[hk] = hv
				}
			}

			alertReq := alertMsgTypeRequst{
				Method: r.Method,
				URL: alertMsgTypeRequstURL{
					Host:     r.URL.Host,
					Path:     r.URL.Path,
					RawQuery: r.URL.RawQuery,
				},
				Host:   r.Host,
				Header: headerMsg,
			}

			alertBaseInfo := alertMsgTypeInfo{
				ServiceName: appinfo.AppInfoIns.Name,
				TraceId:     traceId,
				Id:          appinfo.AppInfoIns.Id,
				Version:     appinfo.AppInfoIns.Version,
			}

			alertMsgObj := alertMsgType{
				Info:        alertBaseInfo,
				HttpRequest: alertReq,
			}

			if contentBuilder != nil {
				alertMsg = contentBuilder(se)
			} else {
				alertMsg = se.Error() + se.Reason + "\n"
				format := "sprintf"

				switch format {
				case "json":
					requestURLJson, urlErr := json.MarshalIndent(r.URL, "", "  ")
					requestHeadersJson, headerErr := json.MarshalIndent(r.Header, "", "  ")
					requestFormJson, formErr := json.MarshalIndent(r.Form, "", "  ")

					if urlErr != nil {
						alertMsg += "request url marshal err:" + urlErr.Error() + "\n"
					}

					if headerErr != nil {
						alertMsg += "request headers marshal err:" + headerErr.Error() + "\n"
					}

					if formErr != nil {
						alertMsg += "request form marshal err:" + formErr.Error() + "\n"
					}

					alertMsg += "request url:\n" + string(requestURLJson) + "\n"
					alertMsg += "request headers:\n" + string(requestHeadersJson) + "\n"
					alertMsg += "request form:\n" + string(requestFormJson) + "\n"

				case "sprintf":

					//formatReqInfo := fmt.Sprintf("%+v", *r)
					reqJsonBytes, marchalErr := json.MarshalIndent(alertMsgObj, "", "  ")
					if marchalErr != nil {
						alertMsg += "alert msg obj marshal err:" + marchalErr.Error() + "\n"
					}

					formatReqInfo := string(reqJsonBytes)

					alertMsg += "\n" + formatReqInfo + "\n"
				}
			}

			//helpers.Alert(alertMsg)

			if alertFun != nil {
				alertFun(alertMsg, alertingIns)
			} else if alertingIns != nil {
				//alertingIns.Alert(alertMsg)

				alertReqMarBytes, marchalErr := json.MarshalIndent(alertReq, "", "  ")
				if marchalErr != nil {
					alertMsg += "alert msg obj marshal err:" + marchalErr.Error() + "\n"
				}

				richMsgMap := types.ChannelRichMsgMap{
					alerting.NotiChanKeyFeishu: feishu.BossAlertCardSvcError(r.Context(), feishu.BossAlertCardVals{
						TraceID: traceId,
						Content: se.Error() + se.Reason,
						Details: string(alertReqMarBytes)}),
					alerting.NotiChanKeyWorkWechat: types.RichMsg{Content: alertMsg},
				}

				alertingIns.AlertRich(richMsgMap)
			} else {
				log.Println("无alerting实例和alertFun")
			}
		}

		w.Write(body)
	}
}

// ContentType returns the content-type with base prefix.
func ContentType(subtype string) string {
	return strings.Join([]string{"application", subtype}, "/")
}
