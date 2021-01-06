/**
 * @Time : 2021/1/6 11:26 上午
 * @Author : MassAdobe
 * @Description: filter
**/
package filter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/errs"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/gin-gonic/gin"
	"github.com/vearne/golib/buffpool"
	"net/http"
	"sync"
	"time"
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 3:15 下午
 * @Description: 自定义返回结构体
**/
type TimeoutWriter struct {
	gin.ResponseWriter
	body        *bytes.Buffer
	h           http.Header
	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 3:15 下午
 * @Description: 修改主返回体
**/
func (tw *TimeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, nil
	}
	return tw.body.Write(b)
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 3:15 下午
 * @Description: 修改返回体头信息
**/
func (tw *TimeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHeader {
		return
	}
	tw.writeHeader(code)
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 3:16 下午
 * @Description: 返回体书写
**/
func (tw *TimeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 3:16 下午
 * @Description: 继承重写方法
**/
func (tw *TimeoutWriter) WriteHeaderNow() {
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 3:16 下午
 * @Description: 继承重写方法
**/
func (tw *TimeoutWriter) Header() http.Header {
	return tw.h
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 3:16 下午
 * @Description: 判断返回体错误代码是否正确
**/
func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(errs.NewError(errs.ErrResponseStatusCode))
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 3:16 下午
 * @Description: 配合gin框架的超时处理
**/
func Timeout(t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		buffer := buffpool.GetBuff()
		c.Writer.Header().Set(constants.CONTENT_TYPE_KEY, constants.CONTENT_TYPE_INNER)
		tw := &TimeoutWriter{body: buffer, ResponseWriter: c.Writer, h: c.Writer.Header()}
		c.Writer = tw
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		c.Request = c.Request.WithContext(ctx)
		finish := make(chan struct{})
		panicChan := make(chan interface{}, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			finish <- struct{}{}
		}()
		select {
		case p := <-panicChan:
			logs.Lg.Error("请求超时", errors.New(fmt.Sprintf("%v", p)), c)
			panic(errs.NewError(errs.ErrSystemCode))
		case <-ctx.Done():
			// 如果超时的话，buffer无法主动清除，只能等待GC回收
			tw.mu.Lock()
			defer tw.mu.Unlock()
			tw.ResponseWriter.WriteHeader(http.StatusRequestTimeout)
			bt, _ := json.Marshal(timeoutError{
				StatusCode: http.StatusRequestTimeout,
				Code:       errs.ErrRequestTimeoutCode,
				Msg:        errs.CodeDescMap[errs.ErrRequestTimeoutCode],
				Data:       nil,
			})
			tw.ResponseWriter.Write(bt)
			c.Abort()
			cancel()
			tw.timedOut = true
			logs.Lg.Error("请求超时", errors.New("request timeout error"), c)
		case <-finish:
			tw.mu.Lock()
			defer tw.mu.Unlock()
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			tw.ResponseWriter.Write(buffer.Bytes())
			buffpool.PutBuff(buffer)
		}
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 3:17 下午
 * @Description: 重新声明返回结构体
**/
type timeoutError struct {
	StatusCode int         `json:"status"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}
