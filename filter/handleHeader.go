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

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-27 21:49
 * @Description: 获取头信息中用户基本信息
**/
func GetReqUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.ToLower(pojo.InitConf.ProgramEnv) != constants.PROGRAM_ENV_DEBUG { // 如果是开发环境
			logs.Lg.SysDebug("中间件-用户信息", c, logs.Desc("当前为非开发环境"))
			if get := c.GetHeader(constants.REQUEST_USER_KEY); len(get) != 0 {
				logs.Lg.SysDebug("中间件-用户信息", c, logs.Desc("头信息中包含用户信息"))
				c.Next()
			} else {
				c.Abort()
				logs.Lg.SysError("中间件-用户信息", errors.New("nil in header"), c, logs.Desc("头信息中不包含用户信息"))
				panic(errs.NewError(errs.ErrHeaderUserCode))
			}
		} else {
			logs.Lg.SysDebug("中间件-用户信息", c, logs.Desc("当前为开发环境"))
			if get := c.GetHeader(constants.REQUEST_USER_KEY); len(get) == 0 {
				logs.Lg.SysDebug("中间件-用户信息", c, logs.Desc("当前开发环境包含头中用户信息"))
				c.Header(constants.REQUEST_USER_KEY, constants.HEADER_USER_JOSN)
			} else {
				logs.Lg.SysDebug("中间件-用户信息", c, logs.Desc("当前开发环境不包含头中用户信息"))
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
			logs.Lg.SysDebug("中间件-跟踪", c, logs.Desc("包含traceId，不用生成"))
		} else {
			traceId := systemUtils.RandomTimestampMark()
			logs.Lg.SysDebug("中间件-跟踪", c, logs.Desc(fmt.Sprintf("不包含traceId，生成Id: %s", traceId)))
			c.Params = append(c.Params, gin.Param{Key: constants.REQUEST_TRACE_ID, Value: traceId})
		}
		// 步骤ID
		if step := c.GetHeader(constants.REQUEST_STEP_ID); 0 != len(step) {
			logs.Lg.SysDebug("中间件-跟踪", c, logs.Desc("包含stepId，不用生成"))
			if parseInt, err := strconv.ParseInt(step, 10, 64); err != nil {
				logs.Lg.SysError("中间件-跟踪", errors.New("marshal stepId error"), c, logs.Desc("解析stepId出错"))
				c.Abort()
			} else {
				stepId := strconv.FormatInt(parseInt+1, 10)
				logs.Lg.SysDebug("中间件-跟踪", c, logs.Desc(fmt.Sprintf("包含stepId，不用生成，自增Id: %s", stepId)))
				c.Params = append(c.Params, gin.Param{Key: constants.REQUEST_STEP_ID, Value: stepId})
			}
		} else {
			logs.Lg.SysDebug("中间件-跟踪", c, logs.Desc(fmt.Sprintf("不包含stepId，生成Id: %s", "0")))
			c.Params = append(c.Params, gin.Param{Key: constants.REQUEST_STEP_ID, Value: "0"})
		}
		// 切面打入调用日志
		if http.MethodPost == c.Request.Method {
			logs.Lg.SysDebug("中间件-跟踪", c, logs.Desc("当前请求为POST请求"))
			if buf, err := ioutil.ReadAll(c.Request.Body); err != nil {
				defer c.Request.Body.Close()
				logs.Lg.SysError("中间件-跟踪", err, c, logs.Desc("关闭请求体失败"))
			} else {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
				logs.Lg.SysInfo("请求日志",
					logs.SpecDesc("请求方法", c.Request.Method),
					logs.SpecDesc("请求路径", c.Request.URL),
					logs.SpecDesc("请求体", fmt.Sprintf("%s", buf)),
					c)
			}
		} else {
			logs.Lg.SysDebug("中间件-跟踪", c, logs.Desc("当前请求为非POST请求"))
			logs.Lg.SysInfo("请求日志",
				logs.SpecDesc("请求方法", c.Request.Method),
				logs.SpecDesc("请求路径", c.Request.URL),
				logs.SpecDesc("请求体", fmt.Sprintf("%s", systemUtils.GetRequestUrlParams(c.Request.RequestURI))),
				c)
		}
	}
}
