/**
 * @Time : 2020/12/17 1:47 下午
 * @Author : MassAdobe
 * @Description: db
**/
package db

import (
	"fmt"
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/nacos"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
	"os"
	"strings"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 1:49 下午
 * @Description: 数据库实体类
**/
var (
	Read  *gorm.DB // 读库
	Write *gorm.DB // 写库
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 1:49 下午
 * @Description: 初始化数据库
**/
func InitDB() {
	if len(nacos.InitConfiguration.Gorm.Read.Ip) != 0 { // 初始化读库
		if gg, err := gorm.Open(constants.DB_TYPE, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			nacos.InitConfiguration.Gorm.Read.Username,
			nacos.InitConfiguration.Gorm.Read.PassWord,
			nacos.InitConfiguration.Gorm.Read.Ip,
			nacos.InitConfiguration.Gorm.Read.Port,
			nacos.InitConfiguration.Gorm.Read.Dbname)); err != nil {
			logs.Lg.SysError("数据库连接", err, logs.Desc("读库连接失败"))
		} else {
			gg.DB().SetMaxIdleConns(2)
			gg.DB().SetMaxOpenConns(10)
			gg.Debug()
			gg.LogMode(true)
			logger := &DbLog{}
			gg.SetLogger(logger)
			if err := gg.DB().Ping(); err != nil {
				logs.Lg.SysError("数据库连接", err, logs.Desc("读库初始化失败"))
				os.Exit(1)
			} else {
				logs.Lg.SysDebug("数据库连接", logs.Desc("读库初始化成功"))
				Read = gg
			}
		}
	}
	if len(nacos.InitConfiguration.Gorm.Write.Ip) != 0 { // 初始化写库
		if gg, err := gorm.Open(constants.DB_TYPE, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			nacos.InitConfiguration.Gorm.Write.Username,
			nacos.InitConfiguration.Gorm.Write.PassWord,
			nacos.InitConfiguration.Gorm.Write.Ip,
			nacos.InitConfiguration.Gorm.Write.Port,
			nacos.InitConfiguration.Gorm.Write.Dbname)); err != nil {
			logs.Lg.SysError("数据库连接", err, logs.Desc("写库连接失败"))
		} else {
			gg.DB().SetMaxIdleConns(2)
			gg.DB().SetMaxOpenConns(10)
			gg.Debug()
			gg.LogMode(true)
			logger := &DbLog{}
			gg.SetLogger(logger)
			if err := gg.DB().Ping(); err != nil {
				logs.Lg.SysError("数据库连接", err, logs.Desc("写库初始化失败"))
				os.Exit(1)
			} else {
				logs.Lg.SysDebug("数据库连接", logs.Desc("写库初始化成功"))
				Write = gg
			}
		}
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/28 2:29 下午
 * @Description: 关停数据库连接池，释放句柄
**/
func CloseDb() {
	mark := false
	if len(nacos.InitConfiguration.Gorm.Read.Ip) != 0 {
		if err := Read.Close(); err != nil {
			mark = true
			logs.Lg.SysError("数据库连接", err, logs.Desc("关闭数据库读库连接池失败"))
		}
	}
	if len(nacos.InitConfiguration.Gorm.Write.Ip) != 0 {
		if err := Write.Close(); err != nil {
			mark = true
			logs.Lg.SysError("数据库连接", err, logs.Desc("关闭数据库写库连接池失败"))
		}
	}
	if mark {
		return
	}
	if len(nacos.InitConfiguration.Gorm.Read.Ip) != 0 || len(nacos.InitConfiguration.Gorm.Write.Ip) != 0 {
		logs.Lg.SysDebug("数据库连接", logs.Desc("关闭数据库连接池成功"))
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/31 1:07 下午
 * @Description: 数据库日志实体类
**/
type DbLog struct{}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/31 1:07 下午
 * @Description: 数据库日志方法
**/
func (logger *DbLog) Print(values ...interface{}) {
	switch values[0] {
	case constants.DB_SQL:
		logs.Lg.SysDebug(
			"数据库日志",
			zap.Any("资源", values[1]),
			zap.String("执行时间", fmt.Sprintf("%v", values[2])),
			zap.String("执行语句", strings.ReplaceAll(values[3].(string), "\n", "")),
			zap.String("执行参数", fmt.Sprintf("%v", values[4])),
			zap.String("影响行数", fmt.Sprintf("%v", values[5])),
		)
	case constants.DB_LOG:
		logs.Lg.SysDebug("数据库日志", zap.Any("其他", values[2]))
	}
}
