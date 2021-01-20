/**
 * @Time : 2021/1/7 10:49 上午
 * @Author : MassAdobe
 * @Description: idempotent
**/
package idempotent

import (
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/errs"
	"github.com/MassAdobe/go-gin/goContext"
	"github.com/MassAdobe/go-gin/rds"
	"github.com/MassAdobe/go-gin/systemUtils"
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/7 10:55 上午
 * @Description: 幂等获取token
**/
func GetToken(c *goContext.Context) {
	token := systemUtils.RandIdempotentToken(c.GetRequestUser().UserId) // 生成幂等的token
	// 获取redis连接
	conn := rds.Get()
	defer conn.Close()
	if _, err := conn.Do(constants.RDS_SETEX, token, constants.TOKEN_EXPIRE_TIME, ""); err != nil {
		c.Error("幂等获取token", err)
		panic(errs.NewError(errs.ErrGetIdempotentCode))
	}
	c.SuccRes(token)
}
