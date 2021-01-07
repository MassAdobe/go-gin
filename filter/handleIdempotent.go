/**
 * @Time : 2021/1/7 2:31 下午
 * @Author : MassAdobe
 * @Description: filter
**/
package filter

import (
	"errors"
	"github.com/MassAdobe/go-gin/errs"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/nacos"
	"github.com/MassAdobe/go-gin/rds"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

const (
	RDS_IDEMPOTENT_HEADER_KEY = "idempotent" // 头中幂等token参数
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/7 2:33 下午
 * @Description: 校验处理幂等接口
**/
func ValidIdempotent() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(nacos.InitConfiguration.Redis.IpPort) != 0 { // 保证幂等 必须存在redis接入
			if key := c.GetHeader(RDS_IDEMPOTENT_HEADER_KEY); len(key) == 0 { // 头中不存在token
				c.Abort()
				logs.Lg.Error("幂等校验", errors.New("header has no token error"), c, logs.Desc("头中不存在相关幂等token"))
				panic(errs.NewError(errs.ErrValidIdempotentHeaderCode))
			} else {
				// 获取redis连接
				conn := rds.Get()
				defer conn.Close()
				// 查询是否存在相关token
				if rply, err := redis.Int(conn.Do(rds.RDS_DEL, key)); err != nil {
					c.Abort()
					logs.Lg.Error("幂等校验", err, c, logs.Desc("redis错误"))
					panic(errs.NewError(errs.ErrValidIdempotentCode))
				} else if 0 == rply { // token不存在，重复提交
					c.Abort()
					logs.Lg.Error("幂等校验", err, c, logs.Desc("redis中不存在相关token"))
					panic(errs.NewError(errs.ErrValidIdempotentRepeatCode))
				} else { // token存在，已经删除，可以继续业务逻辑
					c.Next()
				}
			}
		} else {
			c.Abort()
			logs.Lg.SysError("幂等校验", errors.New("validation idempotent error"), c, logs.Desc("没有接入Redis"))
			panic(errs.NewError(errs.ErrValidIdempotentCode))
		}
	}
}
