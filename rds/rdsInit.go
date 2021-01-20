/**
 * @Time : 2020/12/28 1:40 下午
 * @Author : MassAdobe
 * @Description: rds
**/
package rds

import (
	"errors"
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/nacos"
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

var (
	Redis *redis.Pool
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/28 1:43 下午
 * @Description: 初始化redis连接池
**/
func InitRds() {
	if 0 != len(nacos.InitConfiguration.Redis.IpPort) {
		logs.Lg.SysDebug("Redis连接", logs.Desc("当前nacos配置了redis"))
		if 0 != len(nacos.InitConfiguration.Redis.PassWord) {
			logs.Lg.SysDebug("Redis连接", logs.Desc("当前redis使用密码配置"))
			Redis = &redis.Pool{
				MaxIdle:     nacos.InitConfiguration.Redis.MaxIdle,
				MaxActive:   nacos.InitConfiguration.Redis.MaxActive,
				IdleTimeout: time.Duration(nacos.InitConfiguration.Redis.IdleTimeout) * time.Second,
				Wait:        true,
				Dial: func() (conn redis.Conn, e error) {
					return redis.Dial(constants.REDIS_DIAL_TCP, nacos.InitConfiguration.Redis.IpPort,
						redis.DialPassword(nacos.InitConfiguration.Redis.PassWord),
						redis.DialDatabase(nacos.InitConfiguration.Redis.Database),
						redis.DialConnectTimeout(time.Duration(nacos.InitConfiguration.Redis.ConnectTimeout)*time.Second),
						redis.DialReadTimeout(time.Duration(nacos.InitConfiguration.Redis.ReadTimeout)*time.Second),
						redis.DialWriteTimeout(time.Duration(nacos.InitConfiguration.Redis.WriteTimeout)*time.Second))
				},
			}
			checkConn() // 校验是否成功连接
			return
		}
		logs.Lg.SysDebug("Redis连接", logs.Desc("当前redis使用无密码配置"))
		Redis = &redis.Pool{
			MaxIdle:     nacos.InitConfiguration.Redis.MaxIdle,
			MaxActive:   nacos.InitConfiguration.Redis.MaxActive,
			IdleTimeout: time.Duration(nacos.InitConfiguration.Redis.IdleTimeout) * time.Second,
			Wait:        true,
			Dial: func() (conn redis.Conn, e error) {
				return redis.Dial(constants.REDIS_DIAL_TCP, nacos.InitConfiguration.Redis.IpPort,
					redis.DialDatabase(nacos.InitConfiguration.Redis.Database),
					redis.DialConnectTimeout(time.Duration(nacos.InitConfiguration.Redis.ConnectTimeout)*time.Second),
					redis.DialReadTimeout(time.Duration(nacos.InitConfiguration.Redis.ReadTimeout)*time.Second),
					redis.DialWriteTimeout(time.Duration(nacos.InitConfiguration.Redis.WriteTimeout)*time.Second))
			},
		}
		checkConn() // 校验是否成功连接
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/28 2:15 下午
 * @Description: 校验是否成功连接
**/
func checkConn() {
	logs.Lg.Debug("Redis连接", logs.Desc("检验redis连接是否成功"))
	rc := Get()
	defer rc.Close()
	if reply, err := rc.Do(constants.RDS_PING); err != nil {
		logs.Lg.SysError("Redis连接", err)
		os.Exit(1)
	} else if constants.REDIS_PONG != reply.(string) {
		logs.Lg.SysError("Redis连接", errors.New("redis connect failure"))
		os.Exit(1)
	} else {
		logs.Lg.SysDebug("Redis连接", logs.Desc("Redis连接成功"))
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/28 2:15 下午
 * @Description: 获取redis的连接
**/
func Get() redis.Conn {
	logs.Lg.Debug("Redis连接", logs.Desc("从连接池中获取redis连接"))
	return Redis.Get()
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/28 2:27 下午
 * @Description: 关停redis连接池，释放句柄
**/
func CloseRds() {
	if 0 != len(nacos.InitConfiguration.Redis.IpPort) {
		if err := Redis.Close(); err != nil {
			logs.Lg.SysError("Redis连接", err)
			return
		}
		logs.Lg.SysDebug("Redis连接", logs.Desc("关闭redis连接池成功"))
	}
}
