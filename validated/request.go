/**
 * @Time : 2020-04-30 16:30
 * @Author : MassAdobe
 * @Description: system
**/
package validated

import (
	"github.com/MassAdobe/go-gin/errs"
	"github.com/MassAdobe/go-gin/goContext"
	"github.com/MassAdobe/go-gin/logs"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-30 16:31
 * @Description: 绑定参数并验证参数
**/
func BindAndCheck(c *goContext.Context, data interface{}) {
	if err := c.GinContext.Bind(&data); err != nil { // 获取参数
		logs.Lg.SysError("解析入参", err, c, logs.Desc("解析入参错误"))
		panic(errs.NewError(errs.ErrParamCode))
	}
	logs.Lg.SysDebug("解析入参", c, logs.Desc("解析入参成功"))
	_ = GlobalValidator.Check(data)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-06-01 17:29
 * @Description: 检查参数，一般get请求使用
**/
func CheckParams(data interface{}) {
	_ = GlobalValidator.Check(data)
}
