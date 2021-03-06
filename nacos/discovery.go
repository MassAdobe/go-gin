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
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"os"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 3:16 下午
 * @Description: nacos服务注册发现
**/
func NacosDiscovery() {
	if pojo.InitConf.NacosDiscovery {
		// 创建动态配置客户端
		var namingClientErr error
		// 创建服务发现客户端
		namingClient, namingClientErr = clients.CreateNamingClient(map[string]interface{}{
			constants.NACOS_SERVER_CONFIGS_MARK: serverCs,
			constants.NACOS_CLIENT_CONFIG_MARK:  clientC,
		})
		if nil != namingClientErr {
			logs.Lg.SysError("nacos服务注册与发现", namingClientErr, logs.Desc("创建服务发现客户端失败"))
			os.Exit(1)
		}
		logs.Lg.Debug("nacos服务注册与发现", logs.Desc("创建服务发现客户端成功"))
		if ip, err := systemUtils.ExternalIP(); err != nil {
			logs.Lg.SysError("nacos服务注册与发现", err, logs.Desc("nacos获取当前机器IP失败"))
			os.Exit(1)
		} else {
			pojo.CurIp = ip.String() // 赋值当前宿主IP
			success, namingErr := namingClient.RegisterInstance(vo.RegisterInstanceParam{
				Ip:          pojo.CurIp,
				Port:        InitConfiguration.Serve.Port,
				ServiceName: InitConfiguration.Serve.ServerName,
				Weight:      InitConfiguration.Serve.Weight,
				Enable:      true,
				Healthy:     true,
				Ephemeral:   true,
				Metadata:    map[string]string{constants.NACOS_REGIST_IDC_MARK: constants.NACOS_REGIST_IDC_INNER, constants.NACOS_REGIST_TIMESTAMP_MARK: systemUtils.RtnCurTime(), constants.NACOS_REGIST_VERSION_MARK: pojo.InitConf.CurVersion},
				ClusterName: pojo.InitConf.CurVersion, // 默认值DEFAULT
				GroupName:   pojo.InitConf.NacosGroup, // 默认值DEFAULT_GROUP
			})
			if !success || nil != namingErr {
				logs.Lg.SysError("nacos服务注册与发现", namingErr, logs.Desc("nacos注册服务失败"))
				os.Exit(1)
			}
		}
		logs.Lg.SysDebug("nacos服务注册与发现", logs.Desc("服务注册成功"))
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 4:57 下午
 * @Description: nacos注销服务
**/
func NacosDeregister() {
	if pojo.InitConf.NacosDiscovery {
		success, err := namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          pojo.CurIp,
			Port:        InitConfiguration.Serve.Port,
			ServiceName: InitConfiguration.Serve.ServerName,
			Ephemeral:   true,
			Cluster:     pojo.InitConf.CurVersion, // 默认值DEFAULT
			GroupName:   pojo.InitConf.NacosGroup, // 默认值DEFAULT_GROUP
		})
		if !success || nil != err {
			logs.Lg.SysError("nacos服务注册与发现", err, logs.Desc("nacos注销服务失败"))
			os.Exit(1)
		}
		logs.Lg.SysDebug("nacos服务注册与发现", logs.Desc("nacos服务注销成功"))
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 2:32 下午
 * @Description: 获取服务调用参数
**/
func NacosGetServer(serviceName, groupName string) (instance *model.Instance, err error) {
	instance, err = namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		GroupName:   groupName,                          // 默认值DEFAULT_GROUP
		Clusters:    []string{pojo.InitConf.CurVersion}, // 默认值DEFAULT
	})
	if err != nil {
		logs.Lg.SysError("nacos服务注册与发现", err, logs.Desc(fmt.Sprintf("获取服务失败，查询版本为: %s的服务", pojo.InitConf.CurVersion)))
		if len(pojo.InitConf.LastVersion) != 0 {
			instance, err = namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
				ServiceName: serviceName,
				GroupName:   groupName,                           // 默认值DEFAULT_GROUP
				Clusters:    []string{pojo.InitConf.LastVersion}, // 默认值DEFAULT
			})
			if err != nil {
				logs.Lg.SysError("nacos服务注册与发现", err, logs.Desc(fmt.Sprintf("获取服务再次失败，查询版本为: %s的服务", pojo.InitConf.LastVersion)))
				instance = nil
				return
			}
		} else {
			logs.Lg.SysError("nacos服务注册与发现", errors.New("system has not configure last version"), logs.Desc(fmt.Sprintf("当前系统没有配置上一次版本，查询版本为: %s的服务", pojo.InitConf.LastVersion)))
			instance = nil
			return
		}
	}
	logs.Lg.SysDebug("nacos服务注册与发现", logs.Desc(fmt.Sprintf("获取服务成功: %v", instance)))
	return
}
