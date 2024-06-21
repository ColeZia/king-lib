package log

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"gl.king.im/king-lib/framework/config"
	"gl.king.im/king-lib/framework/internal/di"
	fzap "gl.king.im/king-lib/framework/log/kratoscontrib/zap"
	"gl.king.im/king-lib/framework/log/zlogger"
	"gl.king.im/king-lib/framework/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type Logger struct {
	log.Logger
	channel string
}

var defaultLogger log.Logger
var defaultLogHelper *log.Helper
var esCli *elasticsearch.Client

const esPrefix = "boss-logs-"
const defaultChannel = "service"

// 为 logger 提供写入 redis 队列的 io 接口
type redisLogWriter struct {
	cli     *redis.Client
	listKey string
}

func (w *redisLogWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(w.listKey, p).Result()
	return int(n), err
}

type fileLogWriter struct {
	cli     *redis.Client
	listKey string
}

func (w *fileLogWriter) Write(p []byte) (int, error) {
	n, err := w.cli.RPush(w.listKey, p).Result()
	return int(n), err
}

type esLogWriter struct {
	cli   *elasticsearch.Client
	index string
}

func (w *esLogWriter) Write(p []byte) (int, error) {

	go func(writer *esLogWriter, wp []byte) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				di.Container.DefaultAlerting.Alert(fmt.Sprintln(recoverErr))
			}
		}()

		_, err := writer.writeAsync(wp)
		if err != nil {
			fmt.Println("writer.writeAsync(wp) err:", err)
		}

	}(w, p)

	return len(p), nil
}

func (w *esLogWriter) writeAsync(p []byte) (int, error) {

	// Set up the request object.
	req := esapi.IndexRequest{
		Index: esPrefix + w.index, //"boss-logs",
		//DocumentID: strconv.Itoa(i + 1),
		Body:    strings.NewReader(string(p)),
		Refresh: "true",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := req.Do(ctx, w.cli)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, errors.New(res.String())
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return 0, err
		} else {
			// Print the response status and indexed document version.
			fmt.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	return len(p), err
}

//func (w *esLogWriter) Sync() error {
//	goLog.Println("(w *esLogWriter) Sync():")
//	return nil
//}

var once sync.Once

func GetDefaultLogger() (logger log.Logger, err error) {
	var onceErr error
	once.Do(func() {
		defaultLogger, onceErr = NewLogger(defaultChannel)
	})

	return defaultLogger, onceErr
}

var helperOnce sync.Once

func GetDefaultLogHelper() (h *log.Helper, err error) {
	var onceErr error
	once.Do(func() {
		defaultLogHelper, onceErr = NewLogHelper(defaultChannel)
	})

	return defaultLogHelper, onceErr
}

var esCliOnce sync.Once

func newESCli(conf *config.Bootstrap) (es *elasticsearch.Client, err error) {
	var onceErr error
	once.Do(func() {
		if conf.Service.Log == nil {
			panic("Service.Log配置为空")
		}

		if conf.Service.Log.ES == nil {
			panic("Service.Log.ES配置为空")
		}

		cfg := elasticsearch.Config{
			Addresses: conf.Service.Log.ES.Addrs,
			Username:  conf.Service.Log.ES.Username, // "elastic",
			Password:  conf.Service.Log.ES.Password, //"kGtoDZxJlZnT6YpAKLEr",
			// ...
		}

		fmt.Println("es cfg::", cfg)

		esCli, onceErr = elasticsearch.NewClient(cfg)
		if err != nil {
			panic(err)
		}
	})

	return esCli, onceErr
}

func getServiceConf() *config.Bootstrap {
	//return conf.Bootstrap{Service: &conf.Service{
	//	Log: &conf.Service_Log{
	//		Storage: []conf.Service_Log_STORAGE_TYPE{
	//			conf.Service_Log_STORAGE_ES,
	//			conf.Service_Log_STORAGE_FILE,
	//			conf.Service_Log_STORAGE_STD_OUT,
	//		},
	//		ES: &conf.Service_Log_ES{
	//			Addrs: []string{""},
	//		},
	//	},
	//}}
	return config.GetServiceConf()
}

func NewLogger(channel string) (logger log.Logger, err error) {

	var cores []zapcore.Core

	serviceConf := getServiceConf()
	if serviceConf.Service.Log == nil {
		err = errors.New("Service.Log配置为空")
		return
	}

	if serviceConf.Service.Log.UseZlogger {
		return zlogger.NewLogger(serviceConf.Service.Log)
	}

	// 限制日志输出级别, >= DebugLevel 会打印所有级别的日志
	// 生产环境中一般使用 >= ErrorLevel
	lvCnf := zapcore.InfoLevel

	debugLevelMap := map[config.Service_Log_LogLevel_TYPE]zapcore.Level{
		config.Service_Log_LogLevel_Debug: zapcore.DebugLevel,
		config.Service_Log_LogLevel_Info:  zapcore.InfoLevel,
		config.Service_Log_LogLevel_Warn:  zapcore.WarnLevel,
		config.Service_Log_LogLevel_Error: zapcore.ErrorLevel,
	}

	if val, ok := debugLevelMap[serviceConf.Service.Log.EnableLevel]; ok {
		lvCnf = val
	}

	lowPriority := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= lvCnf
	})

	storateTypeMap := map[config.Service_Log_STORAGE_TYPE]bool{}
	for _, v := range serviceConf.Service.Log.Storage {
		storateTypeMap[v] = true
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.MessageKey = "enc_msg"

	encCfg.TimeKey = "datetime"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	//encCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//	enc.AppendString(t.Format("2006-01-02 15:04:05"))
	//}

	for t, _ := range storateTypeMap {
		switch t {
		case config.Service_Log_STORAGE_STD_OUT:
			// 使用 JSON 格式日志
			jsonEnc := zapcore.NewJSONEncoder(encCfg)

			stdLowPri := lowPriority

			if serviceConf.Service.Log.StdOut != nil {
				if val, ok := debugLevelMap[serviceConf.Service.Log.StdOut.EnableLevel]; ok {
					stdLowPri = zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
						return lv >= val
					})
				}
			}

			stdCore := zapcore.NewCore(jsonEnc, zapcore.Lock(os.Stdout), stdLowPri)

			cores = append(cores, stdCore)

		case config.Service_Log_STORAGE_FILE:
			if serviceConf.Service.Log.File == nil || serviceConf.Service.Log.File.BasePath == "" {
				panic("Service.Log.File.BasePath配置为空")
			}

			//			dir := filepath.Dir(fullFileName)
			//			if _, err := os.Stat(dir); os.IsNotExist(err) {
			//				os.MkdirAll(dir, 0700) // Create your file
			//			}
			//
			//			file, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			//
			//			if err != nil {
			//				log.DefaultLogger.Log(log.LevelError, "title", "fail to open log file!", "err", err)
			//				return nil, err
			//			}

			fullFilePath := filepath.Join(serviceConf.Service.Log.File.BasePath, service.AppInfoIns.Name)

			var writeSyncer zapcore.WriteSyncer
			useRotate := true
			if useRotate {

				sizeBased := true
				if sizeBased {
					fileName := channel + ".log"
					fullFileName := filepath.Join(fullFilePath, fileName)
					//lumberjack.Logger is already safe for concurrent use, so we don't need to lock it.
					writeSyncer = zapcore.AddSync(&lumberjack.Logger{
						Filename: fullFileName,
						MaxSize:  400, // megabytes
						//MaxBackups: 365,
						//MaxAge:     365 * 3, // days
					})
				} else {
					//目前file-rotatelogs这个包已经停止维护-2022-03-27
					fileName := channel
					fullFileName := filepath.Join(fullFilePath, fileName)
					logFile := fullFileName + "-%Y-%m-%d.log"
					rotator, err := rotatelogs.New(
						logFile,
						//rotatelogs.WithRotationCount(365),
						rotatelogs.WithMaxAge(100*365*24*time.Hour),
						rotatelogs.WithRotationTime(time.Hour*24*30),
					)
					if err != nil {
						panic(err)
					}

					writeSyncer = zapcore.AddSync(rotator)
					writeSyncer = zapcore.Lock(writeSyncer)
				}

			} else {
				fullFileName := filepath.Join(fullFilePath, channel+".log")

				var fileClose func()
				writeSyncer, fileClose, err = zap.Open(fullFileName)

				if err != nil {
					fmt.Println("zap.Open::", err)
					fileClose()
					panic(err)
				}

				writeSyncer = zapcore.Lock(writeSyncer)
			}

			jsonEnc := zapcore.NewJSONEncoder(encCfg)
			fileLowPri := lowPriority

			if serviceConf.Service.Log.File != nil {
				if val, ok := debugLevelMap[serviceConf.Service.Log.File.EnableLevel]; ok {
					fileLowPri = zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
						return lv >= val
					})
				}
			}
			stdCore := zapcore.NewCore(jsonEnc, writeSyncer, fileLowPri)

			cores = append(cores, stdCore)

			//logger = log.With(log.NewStdLogger(file),
			//	"ts", log.DefaultTimestamp,
			//	"caller", log.DefaultCaller,
			//	//"service.id", "",
			//	//"service.name", "",
			//	//"service.version", "",
			//	//"trace_id", tracing.TraceID(),
			//	//"span_id", tracing.SpanID(),
			//)
		case config.Service_Log_STORAGE_ES:

			newESCli(getServiceConf())

			esWriter := &esLogWriter{cli: esCli, index: strings.ToLower(service.AppInfoIns.Name)}

			// addSync 将 io.Writer 装饰为 WriteSyncer
			// 故只需要一个实现 io.Writer 接口的对象即可
			syncer := zapcore.AddSync(esWriter)
			jsonEnc := zapcore.NewJSONEncoder(encCfg)
			esLowPri := lowPriority

			if serviceConf.Service.Log.ES != nil {
				if val, ok := debugLevelMap[serviceConf.Service.Log.ES.EnableLevel]; ok {
					esLowPri = zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
						return lv >= val
					})
				}
			}
			esCore := zapcore.NewCore(jsonEnc, syncer, esLowPri)

			cores = append(cores, esCore)

		default:
			panic("不支持的存储类型")
		}
	}

	if len(cores) < 1 {
		panic("未生成日志存储类型")
	}

	// 集成多个 core
	core := zapcore.NewTee(cores...)

	// logger 输出到 console 且标识调用代码行
	lo := zap.New(core).WithOptions(zap.AddCaller(), zap.AddCallerSkip(3))

	lo = lo.With(zap.String("channel", channel))

	logger = fzap.NewLogger(lo)

	return &Logger{
		logger,
		channel,
	}, nil
}

func NewLogHelper(channel string, opts ...log.Option) (h *log.Helper, err error) {
	lo, err := NewLogger(channel)
	if err != nil {
		return nil, err
	}
	h = log.NewHelper(lo)

	return h, nil
}

type ctxLoggerKey struct{}

// NewServerContext creates a new context with client md attached.
func NewLoggerServerContext(ctx context.Context, logger interface{}) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, logger)
}

// FromServerContext returns the server metadata in ctx if it exists.
func CtxLoggerFromServerContext(ctx context.Context) (lo log.Logger, ok bool) {
	lo, ok = ctx.Value(ctxLoggerKey{}).(log.Logger)
	return
}

// FromServerContext returns the server metadata in ctx if it exists.
func LogHelperFromServerContext(ctx context.Context) (loh *log.Helper, ok bool) {
	lo, ok := ctx.Value(ctxLoggerKey{}).(log.Logger)
	if ok {
		loh = log.NewHelper(lo)
	}
	return
}
