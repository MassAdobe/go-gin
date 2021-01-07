/**
 * @Time : 2020/12/17 6:43 下午
 * @Author : MassAdobe
 * @Description: filter
**/
package filter

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/errs"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/pojo"
	"github.com/MassAdobe/go-gin/systemUtils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	HEADER_USER_JOSN  = `{"user_id": 1}` // 如果是开发环境，为了简化开发，默认给定user_id=1
	PROGRAM_ENV_DEBUG = "debug"          // 开发环境
)

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-27 21:49
 * @Description: 获取头信息中用户基本信息
**/
func GetReqUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.ToLower(pojo.InitConf.ProgramEnv) != PROGRAM_ENV_DEBUG { // 如果是开发环境
			if get := c.GetHeader(constants.REQUEST_USER_KEY); len(get) != 0 {
				c.Next()
			} else {
				c.Abort()
				logs.Lg.SysError("头中用户信息为空", errors.New("nil in header"), c)
				panic(errs.NewError(errs.ErrHeaderUserCode))
			}
		} else {
			if get := c.GetHeader(constants.REQUEST_USER_KEY); len(get) == 0 {
				c.Header(constants.REQUEST_USER_KEY, HEADER_USER_JOSN)
			}
			c.Next()
		}
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-27 21:49
 * @Description: 放入请求跟踪ID和步骤ID
**/
func SetTraceAndStep() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跟踪ID
		if trace := c.GetHeader(constants.REQUEST_TRACE_ID); 0 != len(trace) {
			c.Params = append(c.Params, gin.Param{Key: constants.REQUEST_TRACE_ID, Value: trace})
		} else {
			c.Params = append(c.Params, gin.Param{Key: constants.REQUEST_TRACE_ID, Value: systemUtils.RandomTimestampMark()})
		}
		// 步骤ID
		if step := c.GetHeader(constants.REQUEST_STEP_ID); 0 != len(step) {
			if parseInt, err := strconv.ParseInt(step, 10, 64); err != nil {
				logs.Lg.SysError("解析stepId出错", errors.New("Marshal stepId error"))
				c.Abort()
			} else {
				c.Params = append(c.Params, gin.Param{Key: constants.REQUEST_STEP_ID, Value: strconv.FormatInt(parseInt+1, 10)})
			}
		} else {
			c.Params = append(c.Params, gin.Param{Key: constants.REQUEST_STEP_ID, Value: "0"})
		}
		// 切面打入调用日志
		if http.MethodPost == c.Request.Method {
			if buf, err := ioutil.ReadAll(c.Request.Body); err != nil {
				defer c.Request.Body.Close()
			} else {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
				logs.Lg.SysInfo("请求日志",
					logs.SpecDesc("请求方法", c.Request.Method),
					logs.SpecDesc("请求路径", c.Request.URL),
					logs.SpecDesc("请求体", fmt.Sprintf("%s", buf)),
					c)
			}
		} else {
			logs.Lg.SysInfo("请求日志",
				logs.SpecDesc("请求方法", c.Request.Method),
				logs.SpecDesc("请求路径", c.Request.URL),
				logs.SpecDesc("请求体", fmt.Sprintf("%s", systemUtils.GetRequestUrlParams(c.Request.RequestURI))),
				c)
		}
	}
}
