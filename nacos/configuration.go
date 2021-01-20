/**
 * @Time : 2020/12/17 2:25 下午
 * @Author : MassAdobe
 * @Description: nacos
**/
package nacos

import (
	"errors"
	"fmt"
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/pojo"
	"github.com/MassAdobe/go-gin/systemUtils"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"os"
	"strings"
)

var (
	serverCs        []constant.ServerConfig      // nacos的server配置
	clientC         constant.ClientConfig        // nacos的client配置
	profileC        vo.ConfigParam               // nacos的配置
	configClient    config_client.IConfigClient  // nacos服务配置中心client
	namingClient    naming_client.INamingClient  // nacos服务注册与发现client
	NacosContent    string                       // nacos配置中心配置内容
	NacosRegistPojo map[int]interface{}          // nacos注册的自主配置结构体
	RateMap         map[string]ratelimit.Limiter // 限流Map
	PastRateMap     map[string]int               // 历史限流Map
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/8 1:24 下午
 * @Description: 初始化需要初始化的参数
**/
func init() {
	NacosRegistPojo = make(map[int]interface{})  // nacos注册的自主配置结构体
	RateMap = make(map[string]ratelimit.Limiter) // 限流Map
	PastRateMap = make(map[string]int)           // 历史限流Map
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 2:51 下午
 * @Description: 初始化nacos配置
**/
func InitNacos() {
	if pojo.InitConf.NacosConfiguration || pojo.InitConf.NacosDiscovery {
		// 初始化nacos的server服务
		nacosIps := strings.Split(pojo.InitConf.NacosServerIps, ",")
		if 0 == len(nacosIps) {
			fmt.Println(fmt.Sprintf(`{"log_level":"ERROR","time":"%s","msg":"%s","server_name":"%s","desc":"%s"}`, systemUtils.RtnCurTime(), "配置中心", "未知", "nacos地址不能为空"))
			os.Exit(1)
		}
		if 0 == pojo.InitConf.NacosServerPort {
			fmt.Println(fmt.Sprintf(`{"log_level":"ERROR","time":"%s","msg":"%s","server_name":"%s","desc":"%s"}`, systemUtils.RtnCurTime(), "配置中心", "未知", "nacos端口号不能为空"))
			os.Exit(1)
		}
		for _, ip := range nacosIps {
			serverCs = append(serverCs, constant.ServerConfig{
				IpAddr:      ip,
				ContextPath: constants.NACOS_CONTEXT_PATH,
				Port:        pojo.InitConf.NacosServerPort,
				Scheme:      constants.NACOS_SCHEMA,
			})
		}
		// 初始化nacos的client服务
		if 0 == len(pojo.InitConf.NacosClientNamespaceId) {
			fmt.Println(fmt.Sprintf(`{"log_level":"ERROR","time":"%s","msg":"%s","server_name":"%s","desc":"%s"}`, systemUtils.RtnCurTime(), "配置中心", "未知", "nacos命名空间不能为空"))
			os.Exit(1)
		}
		clientC = constant.ClientConfig{
			NamespaceId:         pojo.InitConf.NacosClientNamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
			NotLoadCacheAtStart: constants.NACOS_NOT_LOAD_CACHE_AT_START,
			LogDir:              constants.NACOS_LOG_DIR,
			CacheDir:            constants.NACOS_LOG_CACHE_DIR,
			RotateTime:          constants.NACOS_ROTATE_TIME,
			MaxAge:              constants.NACOS_MAX_AGE,
			LogLevel:            constants.NACOS_LOG_LEVEL,
		}
		if 0 == pojo.InitConf.NacosClientTimeoutMs {
			fmt.Println(fmt.Sprintf(`{"log_level":"ERROR","time":"%s","msg":"%s","server_name":"%s","desc":"%s"}`, systemUtils.RtnCurTime(), "配置中心", "未知", "nacos请求Nacos服务端的超时时间为空，默认为10000ms"))
			os.Exit(1)
		}
		clientC.TimeoutMs = pojo.InitConf.NacosClientTimeoutMs
	}
	if pojo.InitConf.NacosConfiguration {
		// 初始化nacos的获取配置服务
		profileC = vo.ConfigParam{
			DataId: pojo.InitConf.NacosDataId,
			Group:  pojo.InitConf.NacosGroup,
		}
		fmt.Println(fmt.Sprintf(`{"log_level":"INFO","time":"%s","msg":"%s","server_name":"%s","desc":"%s"}`, systemUtils.RtnCurTime(), "配置中心", "未知", "初始化配置成功"))
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 2:50 下午
 * @Description: nacos配置中心
**/
func NacosConfiguration() {
	if pojo.InitConf.NacosConfiguration {
		// 创建动态配置客户端
		var configClientErr error
		configClient, configClientErr = clients.CreateConfigClient(map[string]interface{}{
			constants.NACOS_SERVER_CONFIGS_MARK: serverCs,
			constants.NACOS_CLIENT_CONFIG_MARK:  clientC,
		})
		if nil != configClientErr {
			fmt.Println(fmt.Sprintf(`{"log_level":"ERROR","time":"%s","msg":"%s","server_name":"%s","desc":"%s"}`, systemUtils.RtnCurTime(), "配置中心", "未知", "nacos配置中心连接错误"))
			os.Exit(1)
		}
		// 获取配置
		var contentErr error
		if NacosContent, contentErr = configClient.GetConfig(profileC); contentErr != nil {
			fmt.Println(fmt.Sprintf(`{"log_level":"ERROR","time":"%s","msg":"%s","server_name":"%s","desc":"%s"}`, systemUtils.RtnCurTime(), "配置中心", "未知", "nacos配置中心获取配置错误"))
			os.Exit(1)
		}
		// 返回宿主系统自带参数配置
		for _, v := range NacosRegistPojo {
			ReadNacosSelfProfile(NacosContent, v)
		}
		fmt.Println(fmt.Sprintf(`{"log_level":"INFO","time":"%s","msg":"%s","server_name":"%s","desc":"%s"}`, systemUtils.RtnCurTime(), "配置中心", "未知", "获取配置成功"))
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/21 2:58 下午
 * @Description: 监听配置文件变化
**/
func ListenConfiguration() {
	if pojo.InitConf.NacosConfiguration {
		err := configClient.ListenConfig(vo.ConfigParam{
			DataId: pojo.InitConf.NacosDataId,
			Group:  pojo.InitConf.NacosGroup,
			OnChange: func(namespace, group, dataId, data string) {
				logs.Lg.SysInfo("nacos配置文件监听", logs.Desc(fmt.Sprintf("groupId: %s, dataId: %s, data: %s", group, dataId, data)))
				// 修改日志级别
				profile := ReadNacosProfile(data)
				if strings.ToLower(pojo.InitConf.LogLevel) != strings.ToLower(profile.Log.Level) {
					logs.Lg.SysDebug("nacos配置文件监听", logs.Desc("日志级别修改"))
					switch strings.ToLower(profile.Log.Level) {
					case constants.LOG_LEVEL_MODIFIED_DEBUG:
						logs.Lg.Level.SetLevel(zap.DebugLevel)
						printModifiedLog(profile.Log.Level)
					case constants.LOG_LEVEL_MODIFIED_INFO:
						logs.Lg.Level.SetLevel(zap.InfoLevel)
						printModifiedLog(profile.Log.Level)
					case constants.LOG_LEVEL_MODIFIED_WARN:
						logs.Lg.Level.SetLevel(zap.WarnLevel)
						printModifiedLog(profile.Log.Level)
					case constants.LOG_LEVEL_MODIFIED_ERROR:
						logs.Lg.Level.SetLevel(zap.ErrorLevel)
						printModifiedLog(profile.Log.Level)
					case constants.LOG_LEVEL_MODIFIED_DPANIC:
						logs.Lg.Level.SetLevel(zap.DPanicLevel)
						printModifiedLog(profile.Log.Level)
					case constants.LOG_LEVEL_MODIFIED_PANIC:
						logs.Lg.Level.SetLevel(zap.PanicLevel)
						printModifiedLog(profile.Log.Level)
					case constants.LOG_LEVEL_MODIFIED_FATAL:
						logs.Lg.Level.SetLevel(zap.FatalLevel)
						printModifiedLog(profile.Log.Level)
					default:
						logs.Lg.SysError("动态调整日志级别", errors.New("dynamic modified log level error"), logs.Desc("动态调整日志级别失败，日志级别字符不正确"))
					}
				}
				// 返回宿主系统自带参数配置
				for _, v := range NacosRegistPojo {
					ReadNacosSelfProfile(data, v)
				}
				// 读取限流配置
				ReadRateProfile(profile)
			},
		})
		if err != nil {
			logs.Lg.SysError("nacos配置文件监听", err, logs.Desc("设置nacos配置文件监听器失败"))
			os.Exit(1)
		}
		logs.Lg.SysDebug("nacos配置文件监听", logs.Desc("设置nacos配置文件监听器成功"))
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/21 6:02 下午
 * @Description: 输出动态修改日志级别日志，同时赋值新日志级别
**/
func printModifiedLog(current string) {
	logs.Lg.SysInfo("动态调整日志级别",
		logs.Desc(fmt.Sprintf("动态调整日志级别成功，由级别 %s 调至 %s",
			strings.ToLower(pojo.InitConf.LogLevel), strings.ToLower(current))))
	pojo.InitConf.LogLevel = strings.ToLower(current)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/28 4:07 下午
 * @Description: 插入自定义结构
**/
func InsertSelfProfile(st interface{}) {
	key := len(NacosRegistPojo)
	NacosRegistPojo[key] = st
}
