/**
 * @Time : 2021/1/7 11:55 上午
 * @Author : MassAdobe
 * @Description: routers
**/
package routers

import (
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/errs"
	"github.com/MassAdobe/go-gin/filter"
	"github.com/MassAdobe/go-gin/goContext"
	"github.com/MassAdobe/go-gin/idempotent"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/nacos"
	"github.com/MassAdobe/go-gin/pojo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 10:05 上午
 * @Description: 返回新的router
**/
func NewRouter() (rtr *gin.Engine) {
	logs.Lg.SysDebug("路由", logs.Desc("创建路由"))
	switch strings.ToLower(pojo.InitConf.ProgramEnv) {
	case constants.ROUTER_RELEASE:
		logs.Lg.SysDebug("路由", logs.Desc("当前系统配置启动为生产环境"))
		gin.SetMode(gin.ReleaseMode)
	case constants.ROUTER_TEST:
		logs.Lg.SysDebug("路由", logs.Desc("当前系统配置启动为测试环境"))
		gin.SetMode(gin.TestMode)
	case constants.ROUTER_DEBUG:
		logs.Lg.SysDebug("路由", logs.Desc("当前系统配置启动为开发环境"))
		gin.SetMode(gin.DebugMode)
	default:
		logs.Lg.SysDebug("路由", logs.Desc("当前系统未配置启动环境, 所以默认为生产环境启动"))
		gin.SetMode(gin.ReleaseMode)
	}
	rtr = gin.New()
	rtr.Use(cors.Default()) // 增加跨域处理
	logs.Lg.SysDebug("路由", logs.Desc("增加跨域处理"))
	rtr.NoMethod(errs.HandleNotFound) // 处理没有相关方法时的错误处理
	logs.Lg.SysDebug("路由", logs.Desc("增加处理没有相关方法时的错误处理"))
	rtr.NoRoute(errs.HandleNotFound) // 处理没有相关路由时的错误处理
	logs.Lg.SysDebug("路由", logs.Desc("增加处理没有相关路由时的错误处理"))
	rtr.Use(errs.ErrHandler()) // 全局错误处理
	logs.Lg.SysDebug("路由", logs.Desc("增加全局错误处理"))
	rtr.Use(filter.RateLimit()) // 全局处理限流
	logs.Lg.SysDebug("路由", logs.Desc("增加全局处理限流"))
	other := rtr.Group(nacos.RequestPath("other")).Use(filter.SetTraceAndStep()).Use(filter.GetReqUser())
	{
		// 保证幂等 必须存在redis接入
		if len(nacos.InitConfiguration.Redis.IpPort) != 0 {
			logs.Lg.SysDebug("路由", logs.Desc("增加幂等接口配置功能"))
			other.GET("idempotentToken", goContext.Handle(idempotent.GetToken)) // 幂等获取token
		}
	}
	return rtr
}
