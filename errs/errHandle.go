/**
 * @Time : 2020-04-27 20:14
 * @Author : MassAdobe
 * @Description: error
**/
package errs

import (
	"github.com/gin-gonic/gin"
	"go-gin/logs"
	"net/http"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-27 20:14
 * @Description: 错误处理的结构体
**/
type Error struct {
	StatusCode int         `json:"-"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

var (
	NotFound    = BasicNewError(http.StatusNotFound, ErrNotFoundCode, "", nil)
	ServerError = BasicNewError(http.StatusInternalServerError, ErrSystemCode, "", nil)
)

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-27 20:23
 * @Description: 创建新异常
**/
func NewError(code int, errs ...error) *Error {
	if len(errs) != 0 {
		return BasicNewError(http.StatusOK, code, "", errs[0])
	}
	return BasicNewError(http.StatusOK, code, "", nil)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-30 15:57
 * @Description: 创建存在返回值的新异常
**/
func NewMsgError(code int, msg string, errs ...error) *Error {
	if len(errs) != 0 {
		return BasicNewError(http.StatusOK, code, msg, errs[0])
	}
	return BasicNewError(http.StatusOK, code, msg, nil)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-27 20:25
 * @Description: 其他错误
**/
func OtherError(message string) *Error {
	return &Error{
		StatusCode: http.StatusInternalServerError,
		Code:       ErrSystemCode,
		Msg:        CodeDescMap[ErrSystemCode],
		Data:       message,
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-27 20:20
 * @Description: 基类异常
**/
func BasicNewError(status, code int, msg string, err error) *Error {
	if len(msg) == 0 {
		return &Error{
			StatusCode: status,
			Code:       code,
			Msg:        CodeDescMap[code],
			Data:       "",
		}
	}
	return &Error{
		StatusCode: status,
		Code:       code,
		Msg:        CodeDescMap[code] + "\n" + msg,
		Data:       "",
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-27 20:29
 * @Description: 404处理
**/
func HandleNotFound(c *gin.Context) {
	err := NotFound
	c.JSON(err.StatusCode, err)
	return
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-27 20:29
 * @Description: 整体错误处理
**/
func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var Err *Error
				if e, ok := err.(*Error); ok {
					Err = e
				} else if e, ok := err.(error); ok {
					Err = OtherError(e.Error())
				} else {
					Err = ServerError
				}
				logs.Lg.GlobalError("响应日志", c, logs.Desc("错误"), logs.BasicError(err))
				c.JSON(Err.StatusCode, Err)
				return
			}
		}()
		c.Next()
	}
}