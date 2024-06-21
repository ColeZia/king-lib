package server

import (
	"github.com/gin-gonic/gin"
	ke "github.com/go-kratos/kratos/v2/errors"
	"gl.king.im/king-lib/framework/alerting"
	"gl.king.im/king-lib/framework/alerting/feishu"
	"gl.king.im/king-lib/framework/alerting/types"
	"gl.king.im/king-lib/framework/internal/di"
)

type ResponseStruct struct {
	Code     int         `json:"code"`
	Messsage string      `json:"message"`
	Data     interface{} `json:"data"`
}

func ErrorHandle(ctx *gin.Context, err error, render bool) {
	if err == nil {
		return
	}
	var doAlert bool

	httpStatusCode := 400
	ret := map[string]interface{}{}
	_ = httpStatusCode
	_ = ret

	var errContent string
	switch typeErr := err.(type) {
	case *ke.Error:
		ret = map[string]interface{}{
			"code":     typeErr.Code,
			"reason":   typeErr.Reason,
			"message":  typeErr.Message,
			"metadata": typeErr.Metadata,
		}

		errContent = typeErr.String()

		if typeErr.Code >= 500 && typeErr.Code < 600 {
			doAlert = true
		}

		if typeErr.Code >= 200 && typeErr.Code < 600 {
			httpStatusCode = int(typeErr.Code)
		} else if typeErr.Code >= 600 {
			httpStatusCode = 200
		}

	default:
		ret = map[string]interface{}{
			"code":    400,
			"message": typeErr.Error(),
		}
		errContent = typeErr.Error()
	}

	if doAlert {
		richMsgMap := types.ChannelRichMsgMap{
			alerting.NotiChanKeyFeishu: feishu.BossAlertCardSvcError(ctx.Request.Context(), feishu.BossAlertCardVals{
				Content: errContent,
				Details: ""}),
			alerting.NotiChanKeyWorkWechat: types.RichMsg{Content: errContent},
		}
		di.Container.DefaultAlerting.AlertRich(richMsgMap)
	}

	if render {
		ctx.JSON(httpStatusCode, ret)
	}
}

func RspSerializes(ctx *gin.Context, rspData interface{}, err error) {
	if err == nil {
		ctx.JSON(200, ResponseStruct{
			Code: 200,
			Data: rspData,
		})
	} else {
		ErrorHandle(ctx, err, true)
	}
}

// 以下代码草稿中，请勿使用
type options struct {
	ResponseEncoder ResponseEncoder
}

type Option func(*options)

type ResponseEncoder interface {
	SetData(data interface{})
	SetErrorCode(code int)
	SetReason(reason string)
	SetMessage(msg string)
	SetMetadata(md map[string]string)
	JSON(code int, obj interface{})
	String(code int, format string, values ...interface{})
	HTML(code int, name string, obj interface{})
}

func WithResponseEncoder(re ResponseEncoder) Option {
	return func(o *options) {
		o.ResponseEncoder = re
	}
}

type Serializer interface {
	//RspSerializes(ctx *gin.Context, rspData interface{}, err error)
}

func NewResponseSerializer(inOpts ...Option) Serializer {

	opts := &options{}
	for _, v := range inOpts {
		v(opts)
	}

	return &defaultResponseSerializer{opts: opts}
}

type defaultResponseSerializer struct {
	opts *options
}
