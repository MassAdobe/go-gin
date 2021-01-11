/**
 * @Time : 2021/1/8 1:58 下午
 * @Author : MassAdobe
 * @Description: filter
**/
package filter

import (
	"fmt"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/nacos"
	"github.com/gin-gonic/gin"
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/8 1:58 下午
 * @Description: 处理限流问题
**/
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, okay := nacos.RateMap[nacos.RATE_ALL]; okay { // 如果存在全局
			nacos.RateMap[nacos.RATE_ALL].Take()
			logs.Lg.SysDebug("中间件-限流", logs.Desc("命中全局限流"), c)
		} else if _, okay := nacos.RateMap[c.Request.URL.Path]; okay { // 如果存在当前api
			nacos.RateMap[c.Request.URL.Path].Take()
			logs.Lg.SysDebug("中间件-限流", logs.Desc(fmt.Sprintf("命中%s限流", c.Request.URL.Path)), c)
		} else { // 均无设置
			logs.Lg.SysDebug("中间件-限流", logs.Desc("无限流设置"), c)
			return
		}
	}
}
