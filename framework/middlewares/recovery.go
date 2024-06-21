package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	goLog "log"
	"reflect"
	"runtime/debug"

	kgin "github.com/go-kratos/gin"
	ke "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"gl.king.im/king-lib/framework/alerting"
	"gl.king.im/king-lib/framework/alerting/feishu"
	"gl.king.im/king-lib/framework/alerting/types"
	"gl.king.im/king-lib/framework/internal/di"
	"gl.king.im/king-lib/framework/service"
	"gl.king.im/king-lib/framework/service/appinfo"
)

type recoveryMsgTypeInfo struct {
	ServiceName string
	TraceId     string
	Name        string
	Version     string
	Id          string
	//Env string
}
type recoveryMsgType struct {
	Info recoveryMsgTypeInfo
}

func Recovery() middleware.Middleware {
	return recovery.Recovery(
		recovery.WithLogger(log.DefaultLogger),
		recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
			// do someting

			//helpers.AlertWithTrace(err)

			stack := debug.Stack()
			goLog.Println("Framework Recovery(): err type: ", reflect.TypeOf(err), "; err: ", err, "; stack: ", string(stack))

			var content string

			var errContent string
			switch val := err.(type) {
			case error:
				//return val
				//content = val.Error()

				errContent = val.Error()
			case ke.Error:
				//return &val
				//content = val.Error()
				errContent = val.Error()
			case string:
				//content = fmt.Sprintf("%s", val)
				errContent = fmt.Sprintf("%s", val)
			}

			content += errContent

			traceIdValuer := tracing.TraceID()
			traceId := traceIdValuer(ctx).(string)

			recoveryMsgObj := recoveryMsgType{
				Info: recoveryMsgTypeInfo{
					ServiceName: appinfo.AppInfoIns.Name,
					TraceId:     traceId,
					Id:          appinfo.AppInfoIns.Id,
					Version:     appinfo.AppInfoIns.Version,
				},
			}

			msgBytes, marchalErr := json.MarshalIndent(recoveryMsgObj, "", "  ")
			if marchalErr != nil {
				content += "\n\nrecovery msg obj marshal err:" + marchalErr.Error()
				errContent += "\n\nrecovery msg obj marshal err:" + marchalErr.Error()
			} else {
				content += fmt.Sprintf("\n\n%s", msgBytes)
			}

			content += fmt.Sprintf("\n\nStack::%s", stack)

			if di.Container.DefaultAlerting != nil {
				//di.Container.DefaultAlerting.Alert(content)
				richMsgMap := types.ChannelRichMsgMap{
					alerting.NotiChanKeyFeishu: feishu.BossAlertCardSvcError(ctx, feishu.BossAlertCardVals{
						Content: errContent,
						Details: fmt.Sprintf("\n\nStack::%s", stack)}),
					alerting.NotiChanKeyWorkWechat: types.RichMsg{Content: content},
				}

				di.Container.DefaultAlerting.AlertRich(richMsgMap)
			}

			retErr := ke.New(500, "SERVICE_RECOVERY", service.AppInfoIns.Name+"系统错误！")
			ginCtx, ok := kgin.FromGinContext(ctx)
			if ok {
				ginCtx.Writer.WriteHeader(500)
				ginCtx.Writer.Header().Set("Content-Type", kgin.ContentType("json"))
				kgin.Error(ginCtx, retErr)
			}

			return retErr

		}),
	)
}
