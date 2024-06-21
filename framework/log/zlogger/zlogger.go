package zlogger

import (
	"context"
	"fmt"
	"time"

	"gl.king.im/king-lib/framework/config"

	"git.e.coding.king.cloud/dev/quality/king-micro/zlogger"
	klog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
)

var _ klog.Logger = (*Zlogger)(nil)

type Zlogger struct {
	log *zlogger.LoggerSelector
}

func NewLogger(cnf *config.Service_Log) (logger klog.Logger, err error) {
	//WithExtraCallerSkip(0) 默认为0，可以省略，具体详细看zlogger介绍
	li, err := zlogger.CreateLogger().WithExtraCallerSkip(0).WithLevel(zlogger.LevelInfo.String()).Build()
	if err != nil {
		panic(err)
	}
	//目前所有的stdout都会被采集到logcenter，并自动生成ES索引
	zLog := zlogger.CreateLocalHelper(li) //指定为标准输出
	info := LogInfo{Env: "beta"}
	zLog.Info(nil, info)

	logger = &Zlogger{
		log: zLog,
	}

	return
}

func (l *Zlogger) Log(level klog.Level, keyvals ...interface{}) error {
	ctx := context.Background()
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(ctx, fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case klog.LevelDebug:
		l.log.Debug(ctx, data)
	case klog.LevelInfo:
		l.log.Info(ctx, data)
	case klog.LevelWarn:
		l.log.Warn(ctx, data)
	case klog.LevelError:
		l.log.Error(ctx, data)
	case klog.LevelFatal:
		l.log.Fatal(ctx, data)
	}
	return nil
}

type LogInfo struct {
	LogType      string    `json:"log_type,omitempty"`      //日志类型
	UserId       uint64    `json:"user_id,omitempty"`       //UID-控制台
	Platform     uint32    `json:"platform,omitempty"`      //业务平台
	CallScenario string    `json:"call_scenario,omitempty"` //
	TraceId      string    `json:"trace_id,omitempty"`      //链路ID
	UserName     string    `json:"user_name,omitempty"`     //用户ID
	Duration     float64   `json:"duration,omitempty"`      //时长
	StartTime    time.Time `json:"start_time,omitempty"`    //开始时间
	EndTime      time.Time `json:"end_time,omitempty"`      //结束时间
	Params       string    `json:"params,omitempty"`        //入参
	Output       string    `json:"out_put,omitempty"`       //出参
	Operation    string    `json:"operation,omitempty"`     //接口
	Msg          string    `json:"msg,omitempty"`           //信息
	Code         int64     `json:"code,omitempty"`          //返回码
	Env          string    `json:"env,omitempty"`           //环境
}
