/**
 * @Time : 2020/12/17 4:25 下午
 * @Author : MassAdobe
 * @Description: nacos
**/
package nacos

import "fmt"

var (
	InitConfiguration InitNacosConfiguration // 初始化配置
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 4:26 下午
 * @Description: nacos配置文件配置
**/
type InitNacosConfiguration struct {
	Serve struct { // 服务配置
		Port       uint64  `yaml:"port"`        // 服务端口号
		ServerName string  `yaml:"server-name"` // 服务名
		Weight     float64 `yaml:"weight"`      // nacos中权重
	} `yaml:"serve"`

	Log struct { // 日志配置
		Path  string `yaml:"path"`  // 日志地址
		Level string `yaml:"level"` // 日志级别
	} `yaml:"log"`

	Gorm struct { // 数据库配置
		Read struct { // 读库配置
			Username string `yaml:"username"` // 数据库用户名
			PassWord string `yaml:"password"` // 数据库密码
			Ip       string `yaml:"ip"`       // 数据库IP
			Port     int    `yaml:"port"`     // 数据库端口
			Dbname   string `yaml:"dbname"`   // 数据库名称
		} `yaml:"read"`
		Write struct { // 写库配置
			Username string `yaml:"username"` // 数据库用户名
			PassWord string `yaml:"password"` // 数据库密码
			Ip       string `yaml:"ip"`       // 数据库IP
			Port     int    `yaml:"port"`     // 数据库端口
			Dbname   string `yaml:"dbname"`   // 数据库名称
		} `yaml:"write"`
	} `yaml:"gorm"`
	Feign struct { // 内部调用配置
		RetryNum int `yaml:"retry-num"` // 内部调用重试次数
	} `yaml:"feign"`
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 11:39 上午
 * @Description: 拼装请求主地址
**/
func RequestPath(path string) string {
	return fmt.Sprintf("/%s/%s", InitConfiguration.Serve.ServerName, path)
}
