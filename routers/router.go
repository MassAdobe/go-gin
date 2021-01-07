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
	"github.com/MassAdobe/go-gin/nacos"
	"github.com/MassAdobe/go-gin/pojo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	ROUTER_RELEASE = "release" // 生产
	ROUTER_TEST    = "test"    // 测试
	ROUTER_DEBUG   = "debug"   // 开发
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/6 10:05 上午
 * @Description: 返回新的router
**/
func NewRouter() (rtr *gin.Engine) {
	switch strings.ToLower(pojo.InitConf.ProgramEnv) {
	case ROUTER_RELEASE:
		gin.SetMode(gin.ReleaseMode)
	case ROUTER_TEST:
		gin.SetMode(gin.TestMode)
	case ROUTER_DEBUG:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	rtr = gin.New()
	rtr.Use(cors.Default())           // 增加跨域处理
	rtr.NoMethod(errs.HandleNotFound) // 处理没有相关方法时的错误处理
	rtr.NoRoute(errs.HandleNotFound)  // 处理没有相关路由时的错误处理
	rtr.Use(errs.ErrHandler())        // 全局错误处理
	if gin.Mode() != gin.DebugMode {
		rtr.Use(filter.Timeout(constants.REQUEST_TIMEOUT_TM)) // 增加处理超时请求
	}
	other := rtr.Group("other").Use(filter.SetTraceAndStep()).Use(filter.GetReqUser())
	{
		// 保证幂等 必须存在
		if len(nacos.InitConfiguration.Redis.IpPort) != 0 {
			other.GET("idempotentToken", goContext.Handle(idempotent.GetToken)) // 幂等获取token
		}
	}
	return rtr
}
