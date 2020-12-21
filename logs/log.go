/**
 * @Time : 2020-04-26 19:57
 * @Author : MassAdobe
 * @Description: logs
**/
package logs

import (
	_ "encoding/json"
	"fmt"
	"github.com/MassAdobe/go-gin/config"
	"github.com/MassAdobe/go-gin/pojo"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"reflect"
	"runtime"
	"time"
)

var (
	Lg MyLog
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 6:17 下午
 * @Description: 日志对象
**/
type MyLog struct {
	zapLog *zap.Logger
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 20:02
 * @Description: 日志常量
**/
const (
	TIME              = "TIME"
	LOG_LEVEL         = "LOG-LEVEL"
	LOGGER            = "LOGGER"
	DESC              = "DESC"
	MSG               = "MSG"
	TRACE             = "TRACE"
	ERROR             = "ERROR"
	TIME_FORMAT       = "2006-01-02 15:04:05.000"
	ZAP_FIELD_TYPE    = "zapcore.Field"
	GIN_CONTEXT_TYPE  = "*gin.Context"
	GIN_CONTEXT_TYPES = "[]*gin.Context"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:03
 * @Description: 新建日志
**/
func NewLogger(filePath, level string, maxSize, maxBackups, maxAge int, compress bool, serviceName string) {
	core := newCore(filePath, level, maxSize, maxBackups, maxAge, compress)
	Lg.zapLog = zap.New(core, zap.AddCaller(), zap.Development(), zap.Fields(zap.String("SERVER-NAME", serviceName)))
	zap.ReplaceGlobals(Lg.zapLog)
	Lg.Info("日志启动成功")
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:04
 * @Description: 日志内部配置
**/
func newCore(filePath, level string, maxSize, maxBackups, maxAge int, compress bool) zapcore.Core {
	// 日志文件路径配置
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch level {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
		break
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
		break
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
		break
	default:
		atomicLevel.SetLevel(zap.WarnLevel)
		break
	}
	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        TIME,
		LevelKey:       LOG_LEVEL,
		NameKey:        LOGGER,
		MessageKey:     MSG,
		StacktraceKey:  TRACE,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 颜色编码器
		EncodeTime:     myTimeEncode,                   // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel,                                                                     // 日志级别
	)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:04
 * @Description: 标准化日志日期
**/
func myTimeEncode(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(TIME_FORMAT))
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-27 20:52
 * @Description: 全局处理错误封装
**/
func BasicError(err interface{}) zap.Field {
	return zap.Any(ERROR, err)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:04
 * @Description: 错误封装
**/
func Error(err error) zap.Field {
	if nil == err {
		return zap.String(ERROR, "")
	}
	return zap.String(ERROR, err.Error())
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:04
 * @Description: 描述封装
**/
func Desc(desc interface{}) zap.Field {
	return zap.Any(DESC, desc)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:05
 * @Description: 特殊日志封装
**/
func SpecDesc(name string, desc interface{}) zap.Field {
	return zap.Any(fmt.Sprintf("%s", name), desc)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:05
 * @Description: gin自定义日志插件
**/
func InitLogger(path, serveName, level string, port uint64) {
	if len(pojo.InitConf.LogPath) != 0 && len(pojo.InitConf.LogLevel) != 0 {
		NewLogger(fmt.Sprintf("%s/%s-%d", pojo.InitConf.LogPath, serveName, port),
			pojo.InitConf.LogLevel,
			128,
			10,
			7,
			true,
			serveName)
		return
	}
	NewLogger(fmt.Sprintf("%s/%s-%d", path, serveName, port),
		level,
		128,
		10,
		7,
		true,
		serveName)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 7:52 下午
 * @Description: 重写Global Error日志级别输出
**/
func (this *MyLog) GlobalError(msg string, contextAndFields ...interface{}) {
	fields := this.setTraceAndStep(contextAndFields...)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fields = append(fields, zap.Any("FUNCTION", f.Name()))
	fields = append(fields, zap.Any("PATH-NUM", fmt.Sprintf("%s:%d", file, line)))
	if ce := this.zapLog.Check(zapcore.ErrorLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 7:52 下午
 * @Description: 重写Debug日志级别输出
**/
func (this *MyLog) Debug(msg string, contextAndFields ...interface{}) {
	fields := this.setTraceAndStep(contextAndFields...)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fields = append(fields, zap.Any("FUNCTION", f.Name()))
	fields = append(fields, zap.Any("PATH-NUM", fmt.Sprintf("%s:%d", file, line)))
	if ce := this.zapLog.Check(zapcore.DebugLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 7:52 下午
 * @Description: 重写Info日志级别输出
**/
func (this *MyLog) Info(msg string, contextAndFields ...interface{}) {
	fields := this.setTraceAndStep(contextAndFields...)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fields = append(fields, zap.Any("FUNCTION", f.Name()))
	fields = append(fields, zap.Any("PATH-NUM", fmt.Sprintf("%s:%d", file, line)))
	if ce := this.zapLog.Check(zapcore.InfoLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 7:52 下午
 * @Description: 重写Warn日志级别输出
**/
func (this *MyLog) Warn(msg string, contextAndFields ...interface{}) {
	fields := this.setTraceAndStep(contextAndFields...)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fields = append(fields, zap.Any("FUNCTION", f.Name()))
	fields = append(fields, zap.Any("PATH-NUM", fmt.Sprintf("%s:%d", file, line)))
	if ce := this.zapLog.Check(zapcore.WarnLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 7:52 下午
 * @Description: 重写Error日志级别输出
**/
func (this *MyLog) Error(msg string, err error, contextAndFields ...interface{}) {
	fields := this.setTraceAndStep(contextAndFields...)
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	fields = append(fields, Error(err))
	fields = append(fields, zap.Any("FUNCTION", f.Name()))
	fields = append(fields, zap.Any("PATH-NUM", fmt.Sprintf("%s:%d", file, line)))
	if ce := this.zapLog.Check(zapcore.ErrorLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 7:52 下午
 * @Description: 统一输出trace和step信息
**/
func (this *MyLog) setTraceAndStep(contextAndFields ...interface{}) []zap.Field {
	newFields := make([]zap.Field, 0)
	for _, contextAndField := range contextAndFields {
		t := reflect.TypeOf(contextAndField)
		switch t.String() {
		case GIN_CONTEXT_TYPE:
			context := contextAndField.(*gin.Context)
			if value, has := context.Params.Get(config.REQUEST_TRACE_ID); has {
				newFields = append(newFields, zap.Any(config.REQUEST_TRACE_ID, value))
			}
			if value, has := context.Params.Get(config.REQUEST_STEP_ID); has {
				newFields = append(newFields, zap.Any(config.REQUEST_STEP_ID, value))
			}
		case GIN_CONTEXT_TYPES:
			context := contextAndField.([]*gin.Context)[0]
			if value, has := context.Params.Get(config.REQUEST_TRACE_ID); has {
				newFields = append(newFields, zap.Any(config.REQUEST_TRACE_ID, value))
			}
			if value, has := context.Params.Get(config.REQUEST_STEP_ID); has {
				newFields = append(newFields, zap.Any(config.REQUEST_STEP_ID, value))
			}
		case ZAP_FIELD_TYPE:
			newFields = append(newFields, contextAndField.(zapcore.Field))
		default:
		}
	}
	return newFields
}
