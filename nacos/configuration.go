/**
 * @Time : 2020/12/17 2:25 下午
 * @Author : MassAdobe
 * @Description: nacos
**/
package nacos

import (
	"fmt"
	"github.com/MassAdobe/go-gin/pojo"
	"github.com/MassAdobe/go-gin/systemUtils"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"os"
	"strings"
)

var (
	serverCs     []constant.ServerConfig     // nacos的server配置
	clientC      constant.ClientConfig       // nacos的client配置
	profileC     vo.ConfigParam              // nacos的配置
	configClient config_client.IConfigClient // nacos服务配置中心client
	namingClient naming_client.INamingClient // nacos服务注册与发现client
	NacosContent string                      // nacos配置中心配置内容
)

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
			fmt.Println(fmt.Sprintf("【SYSTEM】%s %s %s %s", systemUtils.RtnCurTime(), "【nacos配置中心】", "【nacos配置错误】", "nacos地址不能为空"))
			os.Exit(1)
		}
		if 0 == pojo.InitConf.NacosServerPort {
			fmt.Println(fmt.Sprintf("【SYSTEM】%s %s %s %s", systemUtils.RtnCurTime(), "【nacos配置中心】", "【nacos配置错误】", "nacos端口号不能为空"))
			os.Exit(1)
		}
		for _, ip := range nacosIps {
			serverCs = append(serverCs, constant.ServerConfig{
				IpAddr:      ip,
				ContextPath: "/nacos",
				Port:        pojo.InitConf.NacosServerPort,
				Scheme:      "http",
			})
		}
		// 初始化nacos的client服务
		if 0 == len(pojo.InitConf.NacosClientNamespaceId) {
			fmt.Println(fmt.Sprintf("【SYSTEM】%s %s %s %s", systemUtils.RtnCurTime(), "【nacos配置中心】", "【nacos配置错误】", "nacos命名空间不能为空"))
			os.Exit(1)
		}
		clientC = constant.ClientConfig{
			NamespaceId:         pojo.InitConf.NacosClientNamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
			NotLoadCacheAtStart: true,
			LogDir:              "/tmp/nacos/log",
			CacheDir:            "/tmp/nacos/cache",
			RotateTime:          "1h",
			MaxAge:              3,
			LogLevel:            "debug",
		}
		if 0 == pojo.InitConf.NacosClientTimeoutMs {
			fmt.Println(fmt.Sprintf("【SYSTEM】%s %s %s", systemUtils.RtnCurTime(), "【nacos配置中心】", "nacos请求Nacos服务端的超时时间为空，默认为10000ms"))
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
	}
	fmt.Println(fmt.Sprintf("【SYSTEM】%s %s %s %s", systemUtils.RtnCurTime(), "【nacos配置中心】", "【nacos配置】", "初始化配置成功"))
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
			"serverConfigs": serverCs,
			"clientConfig":  clientC,
		})
		if nil != configClientErr {
			fmt.Println(fmt.Sprintf("【SYSTEM】%s %s %v %s", systemUtils.RtnCurTime(), "【nacos配置中心】", configClientErr, "nacos配置中心连接错误"))
			os.Exit(1)
		}
		// 获取配置
		var contentErr error
		if NacosContent, contentErr = configClient.GetConfig(profileC); contentErr != nil {
			fmt.Println(fmt.Sprintf("【SYSTEM】%s %s %v %s", systemUtils.RtnCurTime(), "【nacos配置中心】", contentErr, "nacos配置中心获取配置错误"))
			os.Exit(1)
		}
		fmt.Println(fmt.Sprintf("【SYSTEM】%s %s %s %s", systemUtils.RtnCurTime(), "【nacos配置中心】", "【nacos配置】", "获取配置成功"))
	}
}
