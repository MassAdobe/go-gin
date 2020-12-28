/**
 * @Time : 2020/12/28 1:40 下午
 * @Author : MassAdobe
 * @Description: rds
**/
package rds

import (
	"errors"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/nacos"
	"github.com/gomodule/redigo/redis"
	"os"
	"time"
)

const (
	REDIS_DIAL_TCP = "tcp"
)

var (
	Redis *redis.Pool
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/28 1:41 下午
 * @Description: redis的命令常量
**/
const (
	RDS_PING     = "PING"
	RDS_HEXISTS  = "hexists"
	RDS_HSET     = "hset"
	RDS_HGET     = "hget"
	RDS_HDEL     = "hdel"
	RDS_HGETALL  = "hgetall"
	RDS_HKEYS    = "hkeys"
	RDS_EXISTS   = "exists"
	RDS_GET      = "get"
	RDS_SET      = "set"
	RDS_ZADD     = "zadd"
	RDS_ZREVRANK = "zrevrank"
	RDS_ZCOUNT   = "zcount"
	RDS_DEL      = "del"
	RDS_SETEX    = "setex"
	RDS_HLEN     = "hlen"
	RDS_INCR     = "incr"
	RDS_HINCRBY  = "hincrby"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/28 1:43 下午
 * @Description: 初始化redis连接池
**/
func InitRds() {
	if 0 != len(nacos.InitConfiguration.Redis.IpPort) {
		if 0 != len(nacos.InitConfiguration.Redis.PassWord) {
			Redis = &redis.Pool{
				MaxIdle:     nacos.InitConfiguration.Redis.MaxIdle,
				MaxActive:   nacos.InitConfiguration.Redis.MaxActive,
				IdleTimeout: time.Duration(nacos.InitConfiguration.Redis.IdleTimeout) * time.Second,
				Wait:        true,
				Dial: func() (conn redis.Conn, e error) {
					return redis.Dial(REDIS_DIAL_TCP, nacos.InitConfiguration.Redis.IpPort,
						redis.DialPassword(nacos.InitConfiguration.Redis.PassWord),
						redis.DialDatabase(nacos.InitConfiguration.Redis.Database),
						redis.DialConnectTimeout(time.Duration(nacos.InitConfiguration.Redis.ConnectTimeout)*time.Second),
						redis.DialReadTimeout(time.Duration(nacos.InitConfiguration.Redis.ReadTimeout)*time.Second),
						redis.DialWriteTimeout(time.Duration(nacos.InitConfiguration.Redis.WriteTimeout)*time.Second))
				},
			}
			return
		}
		Redis = &redis.Pool{
			MaxIdle:     nacos.InitConfiguration.Redis.MaxIdle,
			MaxActive:   nacos.InitConfiguration.Redis.MaxActive,
			IdleTimeout: time.Duration(nacos.InitConfiguration.Redis.IdleTimeout) * time.Second,
			Wait:        true,
			Dial: func() (conn redis.Conn, e error) {
				return redis.Dial(REDIS_DIAL_TCP, nacos.InitConfiguration.Redis.IpPort,
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
	rc := Get()
	defer rc.Close()
	if reply, err := rc.Do(RDS_PING); err != nil {
		logs.Lg.Error("Redis连接", err)
		os.Exit(1)
	} else if "PONG" != reply.(string) {
		logs.Lg.Error("Redis连接", errors.New("Redis连接校验失败"))
		os.Exit(1)
	} else {
		logs.Lg.Info("Redis连接", logs.Desc("Redis连接成功"))
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/28 2:15 下午
 * @Description: 获取redis的连接
**/
func Get() redis.Conn {
	return Redis.Get()
}
