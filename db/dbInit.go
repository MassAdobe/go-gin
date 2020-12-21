/**
 * @Time : 2020/12/17 1:47 下午
 * @Author : MassAdobe
 * @Description: db
**/
package db

import (
	"fmt"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/nacos"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
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
		if gg, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			nacos.InitConfiguration.Gorm.Read.Username,
			nacos.InitConfiguration.Gorm.Read.PassWord,
			nacos.InitConfiguration.Gorm.Read.Ip,
			nacos.InitConfiguration.Gorm.Read.Port,
			nacos.InitConfiguration.Gorm.Read.Dbname)); err != nil {
			logs.Lg.Error("数据库读库连接失败", err)
		} else {
			gg.DB().SetMaxIdleConns(2)
			gg.DB().SetMaxOpenConns(10)
			gg.LogMode(true)
			if err := gg.DB().Ping(); err != nil {
				logs.Lg.Error("数据库读库初始化失败", err)
				os.Exit(1)
			} else {
				logs.Lg.Info("数据库读库初始化成功")
				Read = gg
			}
		}
	}
	if len(nacos.InitConfiguration.Gorm.Write.Ip) != 0 { // 初始化写库
		if gg, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			nacos.InitConfiguration.Gorm.Write.Username,
			nacos.InitConfiguration.Gorm.Write.PassWord,
			nacos.InitConfiguration.Gorm.Write.Ip,
			nacos.InitConfiguration.Gorm.Write.Port,
			nacos.InitConfiguration.Gorm.Write.Dbname)); err != nil {
			logs.Lg.Error("数据库写库连接失败", err)
		} else {
			gg.DB().SetMaxIdleConns(2)
			gg.DB().SetMaxOpenConns(10)
			gg.LogMode(true)
			if err := gg.DB().Ping(); err != nil {
				logs.Lg.Error("数据库写库初始化失败", err)
				os.Exit(1)
			} else {
				logs.Lg.Info("数据库写库初始化成功")
				Write = gg
			}
		}
	}
}
