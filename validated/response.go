/**
 * @Time : 2020-04-26 19:57
 * @Author : MassAdobe
 * @Description: system
**/
package validated

import (
	"github.com/MassAdobe/go-gin/errs"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
func SuccRes(c *gin.Context, data interface{}) {
	logs.Lg.Info("响应日志",
		logs.SpecDesc("请求方法", c.Request.Method),
		logs.SpecDesc("请求路径", c.Request.URL),
		logs.SpecDesc("响应体", data),
		c)
	c.JSON(http.StatusOK, res(errs.SuccessCode, data))
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:12
 * @Description: 错误时返回 支持data为空
**/
func FailRes(c *gin.Context, errCode int, data interface{}) {
	logs.Lg.Info("响应日志",
		logs.SpecDesc("请求方法", c.Request.Method),
		logs.SpecDesc("请求路径", c.Request.URL),
		logs.SpecDesc("响应体", data),
		c)
	c.JSON(http.StatusOK, res(errCode, data))
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 6:19 下午
 * @Description: 内部调用 成功时返回 支持data为空
**/
func SuccResFeign(c *gin.Context, data interface{}) {
	logs.Lg.Info("响应日志",
		logs.SpecDesc("请求方法", c.Request.Method),
		logs.SpecDesc("请求路径", c.Request.URL),
		logs.SpecDesc("响应体", data),
		c)
	c.JSON(http.StatusOK, data)
}
