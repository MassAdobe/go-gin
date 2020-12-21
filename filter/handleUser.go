/**
 * @Time : 2020/12/18 3:24 下午
 * @Author : MassAdobe
 * @Description: filter
**/
package filter

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-gin/config"
	"go-gin/errs"
	"go-gin/logs"
	"go-gin/pojo"
	"net/url"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-28 10:21
 * @Description: 获取用户基本信息
**/
func GetRequestUser(c *gin.Context) *pojo.RequestUser {
	rq := c.GetHeader(config.REQUEST_USER_KEY)
	if len(rq) != 0 {
		c.Params = append(c.Params, gin.Param{Key: config.REQUEST_USER_KEY, Value: rq})
		if enEscapeUrl, err := url.QueryUnescape(rq); err != nil {
			logs.Lg.Error("解析头中用户信息错误", err, c)
			panic(errs.NewError(errs.ErrHeaderUserCode, err))
		} else {
			requestUser := new(pojo.RequestUser)
			if err := json.Unmarshal([]byte(enEscapeUrl), &requestUser); err != nil {
				logs.Lg.Error("解析头中用户信息JSON错误", err, c)
				panic(errs.NewError(errs.ErrHeaderUserCode, err))
			}
			return requestUser
		}
	}
	return nil
}
