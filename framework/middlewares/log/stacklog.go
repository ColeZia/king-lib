package middlewares

import (
	"context"

	ke "github.com/go-kratos/kratos/v2/errors"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"gl.king.im/king-lib/framework/config"
	"gl.king.im/king-lib/framework/internal/di"
	"gl.king.im/king-lib/framework/log"
)

func mapLogLevel(ll config.Service_Log_LogLevel_TYPE) klog.Level {
	switch ll {
	case config.Service_Log_LogLevel_Unknown:
		return klog.LevelInfo
	case config.Service_Log_LogLevel_Debug:
		return klog.LevelDebug
	case config.Service_Log_LogLevel_Info:
		return klog.LevelInfo
	case config.Service_Log_LogLevel_Warn:
		return klog.LevelWarn
	case config.Service_Log_LogLevel_Error:
		return klog.LevelError
	case config.Service_Log_LogLevel_Panic:
		return klog.LevelError
	case config.Service_Log_LogLevel_Fatal:
		return klog.LevelFatal
	}
	return klog.LevelInfo
}

func StackLog() middleware.Middleware {
	svcCnf := config.GetServiceConf()
	logLevel := klog.LevelInfo
	if svcCnf.Service.Log != nil && svcCnf.Service.Log.EnableLevel != config.Service_Log_LogLevel_Unknown {
		logLevel = mapLogLevel(svcCnf.Service.Log.EnableLevel)
	}

	return func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			logger := di.GetLogger()
			l2 := klog.WithContext(ctx, logger)
			l2 = klog.With(l2, "module", "framework/middleware/log/StackLog")

			stackLog := log.NewStackLog(logLevel)
			stackLog.AppendLog(klog.LevelInfo, "step", "begin")
			ctx = log.NewStackLogServerContext(ctx, stackLog)

			rsp, err := next(ctx, req)

			var kvs []interface{}

			retCode := 200
			if err != nil {
				stackLog.AppendLog(klog.LevelInfo, "error_msg", err.Error())
				switch val := err.(type) {
				case *ke.Error:
					retCode = int(val.Code)
				case error:
					retCode = 400
				}
			}

			stackLog.AppendLog(klog.LevelInfo, "step", "end")

			kvs = append(kvs, "ret_code", retCode)

			if len(kvs) > 1 {
				l2 = klog.With(l2, kvs...)
			}

			customModule := stackLog.GetCustomModule()
			if customModule != "" {
				l2 = klog.With(l2, "module", customModule)
			}

			l2.Log(klog.LevelInfo, "stack_log", stackLog.GetLogs())

			return rsp, err
		}
	}
}
