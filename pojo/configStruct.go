/**
 * @Time : 2020-04-26 19:57
 * @Author : MassAdobe
 * @Description: pojo
**/
package pojo

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 2:06 下午
 * @Description: 配置项目
**/
var (
	InitConf InitConfig // 初始化配置
	SysConf  SysConfig  // 系统配置
	CurIp    string     // 当前宿主IP
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 2:06 下午
 * @Description: 初始化配置
**/
type InitConfig struct {
	NacosConfiguration     bool   `yaml:"NacosConfiguration"`     // 是否开启nacos配置中心
	NacosDiscovery         bool   `yaml:"NacosDiscovery"`         // 是否开启nacos服务注册于发现
	NacosServerIps         string `yaml:"NacosServerIps"`         // nacos地址
	NacosServerPort        uint64 `yaml:"NacosServerPort"`        // nacos端口号
	NacosClientNamespaceId string `yaml:"NacosClientNamespaceId"` // nacos命名空间
	NacosClientTimeoutMs   uint64 `yaml:"NacosClientTimeoutMs"`   // 请求Nacos服务端的超时时间，默认是10000ms
	NacosDataId            string `yaml:"NacosDataId"`            // nacos配置文件名称
	NacosGroup             string `yaml:"NacosGroup"`             // nacos配置组名称
	LogPath                string `yaml:"LogPath"`                // 日志输出路径(本地配置优先级最高)
	LogLevel               string `yaml:"LogLevel"`               // 日志级别(本地配置优先级最高)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 19:59
 * @Description: 系统配置
**/
type SysConfig struct {
	LogPath             string `yaml:"LogPath"`             // 日志地址
	LogLevel            string `yaml:"LogLevel"`            // 日志级别
	RegisterIp          string `yaml:"RegisterIp"`          // 注册IP地址
	RedisMaxIdle        int    `yaml:"RedisMaxIdle"`        // Redis最大挂起数
	RedisMaxActive      int    `yaml:"RedisMaxActive"`      // Redis最大活跃数
	RedisMaxIdleTimeout int    `yaml:"RedisMaxIdleTimeout"` // Redis最大挂起时间
	RedisHost           string `yaml:"RedisHost"`           // Redis地址
	RedisPassword       string `yaml:"RedisPassword"`       // Redis密码
	RedisDb             int    `yaml:"RedisDb"`             // Redis数据库
	RedisConnectTimeout int    `yaml:"RedisConnectTimeout"` // Redis连接超时
	RedisReadTimeout    int    `yaml:"RedisReadTimeout"`    // Redis读超时
	RedisWriteTimeout   int    `yaml:"RedisWriteTimeout"`   // Redis写超时
	TokenVerify         string `yaml:"TokenVerify"`         // Token中校验元素secret
	JwtKey              string `yaml:"JwtKey"`              // JWT认证加密私钥
	MysqlUser           string `yaml:"MysqlUser"`           // 数据库用户名
	MysqlPassword       string `yaml:"MysqlPassword"`       // 数据库密码
	MysqlHost           string `yaml:"MysqlHost"`           // 数据库IP
	MysqlPort           string `yaml:"MysqlPort"`           // 数据库端口
	MysqlDbName         string `yaml:"MysqlDbName"`         // 数据库名称
}
