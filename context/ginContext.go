/**
 * @Time : 2021/1/5 5:09 下午
 * @Author : MassAdobe
 * @Description: context
**/
package context

import (
	"fmt"
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/errs"
	"github.com/MassAdobe/go-gin/filter"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/pojo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
	"strings"
)

const (
	ROUTER_RELEASE = "release"
	ROUTER_TEST    = "test"
	ROUTER_DEBUG   = "debug"
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 10:05 上午
 * @Description: 返回新的router
**/
func NewRouter() (rtr *gin.Engine) {
	switch strings.ToLower(pojo.InitConf.ProgramEnv) {
	case ROUTER_RELEASE:
		gin.SetMode(gin.ReleaseMode)
	case ROUTER_TEST:
		gin.SetMode(gin.TestMode)
	case ROUTER_DEBUG:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	rtr = gin.New()
	rtr.NoMethod(errs.HandleNotFound) // 处理没有相关方法时的错误处理
	rtr.NoRoute(errs.HandleNotFound)  // 处理没有相关路由时的错误处理
	rtr.Use(errs.ErrHandler())        // 全局错误处理
	if gin.Mode() != gin.DebugMode {
		rtr.Use(filter.Timeout(constants.REQUEST_TIMEOUT_TM)) // 增加处理超时请求
	}
	return rtr
}

type Context struct {
	GinContext *gin.Context
	GinLog     *logs.MyLog
}

type HandlerFunc func(c *Context)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 10:03 上午
 * @Description: 处理日志与gin框架合并
**/
func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{c, &logs.Lg}
		h(ctx)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 7:52 下午
 * @Description: 重写Debug日志级别输出
**/
func (this *Context) Debug(msg string, fields ...zap.Field) {
	newFields := this.setTraceAndStep()
	if len(fields) > 0 {
		newFields = append(newFields, fields...)
	}
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	newFields = append(newFields, zap.Any("function", f.Name()))
	newFields = append(newFields, zap.Any("path_num", fmt.Sprintf("%s:%d", file, line)))
	if ce := this.GinLog.ZapLog.Check(zapcore.DebugLevel, msg); ce != nil {
		ce.Write(newFields...)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 7:52 下午
 * @Description: 重写Info日志级别输出
**/
func (this *Context) Info(msg string, fields ...zap.Field) {
	newFields := this.setTraceAndStep()
	if len(fields) > 0 {
		newFields = append(newFields, fields...)
	}
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	newFields = append(newFields, zap.Any("function", f.Name()))
	newFields = append(newFields, zap.Any("path_num", fmt.Sprintf("%s:%d", file, line)))
	if ce := this.GinLog.ZapLog.Check(zapcore.InfoLevel, msg); ce != nil {
		ce.Write(newFields...)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 7:52 下午
 * @Description: 重写Warn日志级别输出
**/
func (this *Context) Warn(msg string, fields ...zap.Field) {
	newFields := this.setTraceAndStep()
	if len(fields) > 0 {
		newFields = append(newFields, fields...)
	}
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	newFields = append(newFields, zap.Any("function", f.Name()))
	newFields = append(newFields, zap.Any("path_num", fmt.Sprintf("%s:%d", file, line)))
	if ce := this.GinLog.ZapLog.Check(zapcore.WarnLevel, msg); ce != nil {
		ce.Write(newFields...)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 7:52 下午
 * @Description: 重写Error日志级别输出
**/
func (this *Context) Error(msg string, err error, fields ...zap.Field) {
	newFields := this.setTraceAndStep()
	if len(fields) > 0 {
		newFields = append(newFields, fields...)
	}
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	newFields = append(newFields, logs.Error(err))
	newFields = append(newFields, zap.Any("function", f.Name()))
	newFields = append(newFields, zap.Any("path_num", fmt.Sprintf("%s:%d", file, line)))
	if ce := this.GinLog.ZapLog.Check(zapcore.ErrorLevel, msg); ce != nil {
		ce.Write(newFields...)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 7:52 下午
 * @Description: 统一输出trace和step信息
**/
func (this *Context) setTraceAndStep() []zap.Field {
	newFields := make([]zap.Field, 0)
	if value, has := this.GinContext.Params.Get(constants.REQUEST_TRACE_ID); has {
		newFields = append(newFields, zap.Any(constants.REQUEST_TRACE_ID, value))
	}
	if value, has := this.GinContext.Params.Get(constants.REQUEST_STEP_ID); has {
		newFields = append(newFields, zap.Any(constants.REQUEST_STEP_ID, value))
	}
	return newFields
}
