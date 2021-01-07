/**
 * @Time : 2021/1/5 5:09 下午
 * @Author : MassAdobe
 * @Description: context
**/
package goContext

import (
	"encoding/json"
	"fmt"
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/errs"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/pojo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/url"
	"runtime"
)

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
func Handle(handle HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handle(&Context{c, &logs.Lg})
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-28 10:21
 * @Description: 获取用户基本信息
**/
func (this *Context) GetRequestUser() *pojo.RequestUser {
	rq := this.GinContext.GetHeader(constants.REQUEST_USER_KEY)
	if len(rq) != 0 {
		if enEscapeUrl, err := url.QueryUnescape(rq); err != nil {
			this.SysError("解析头中用户信息错误", err)
			panic(errs.NewError(errs.ErrHeaderUserCode, err))
		} else {
			requestUser := new(pojo.RequestUser)
			if err := json.Unmarshal([]byte(enEscapeUrl), &requestUser); err != nil {
				this.SysError("解析头中用户信息JSON错误", err)
				panic(errs.NewError(errs.ErrHeaderUserCode, err))
			}
			return requestUser
		}
	}
	return nil
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
 * @TIME: 2021/1/7 11:59 上午
 * @Description: 重写Info日志级别输出（系统）
**/
func (this *Context) SysInfo(msg string, fields ...zap.Field) {
	newFields := this.setTraceAndStep()
	if len(fields) > 0 {
		newFields = append(newFields, fields...)
	}
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
 * @Description: 重写Error日志级别输出（系统）
**/
func (this *Context) SysError(msg string, err error, fields ...zap.Field) {
	newFields := this.setTraceAndStep()
	if len(fields) > 0 {
		newFields = append(newFields, fields...)
	}
	newFields = append(newFields, logs.Error(err))
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

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:12
 * @Description: 通用返回结构体封账
**/
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:12
 * @Description: 抽象方法 返回结构体
**/
func res(code int, data interface{}) (rtn *Response) {
	if nil != data {
		rtn = &Response{code, errs.CodeDescMap[code], data}
	} else {
		rtn = &Response{code, errs.CodeDescMap[code], ""}
	}
	return
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:12
 * @Description: 成功时返回 支持data为空
**/
func (this *Context) SuccRes(data interface{}) {
	this.SysInfo("响应日志",
		logs.SpecDesc("请求方法", this.GinContext.Request.Method),
		logs.SpecDesc("请求路径", this.GinContext.Request.URL),
		logs.SpecDesc("响应体", data))
	this.GinContext.JSON(http.StatusOK, res(errs.SuccessCode, data))
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:12
 * @Description: 错误时返回 支持data为空
**/
func (this *Context) FailRes(errCode int, data interface{}) {
	this.SysInfo("响应日志",
		logs.SpecDesc("请求方法", this.GinContext.Request.Method),
		logs.SpecDesc("请求路径", this.GinContext.Request.URL),
		logs.SpecDesc("响应体", data))
	this.GinContext.JSON(http.StatusOK, res(errCode, data))
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 6:19 下午
 * @Description: 内部调用 成功时返回 支持data为空
**/
func (this *Context) SuccResFeign(data interface{}) {
	this.SysInfo("响应日志",
		logs.SpecDesc("请求方法", this.GinContext.Request.Method),
		logs.SpecDesc("请求路径", this.GinContext.Request.URL),
		logs.SpecDesc("响应体", data))
	this.GinContext.JSON(http.StatusOK, data)
}
