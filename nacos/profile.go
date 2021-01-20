/**
 * @Time : 2020/12/17 4:24 下午
 * @Author : MassAdobe
 * @Description: nacos
**/
package nacos

import (
	"fmt"
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/systemUtils"
	"go.uber.org/ratelimit"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 4:24 下午
 * @Description: 处理首次nacos获取到的配置信息
**/
func InitNacosProfile() {
	if err := yaml.Unmarshal([]byte(NacosContent), &InitConfiguration); err != nil {
		fmt.Println(fmt.Sprintf(`{"log_level":"ERROR","time":"%s","msg":"%s","server_name":"%s","desc":"%s"}`, systemUtils.RtnCurTime(), "配置中心", "未知", "读取nacos系统配置失败"))
		os.Exit(1)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/21 3:05 下午
 * @Description: 返回配置文件内容
**/
func ReadNacosProfile(content string) *InitNacosConfiguration {
	var NewInitConfiguration InitNacosConfiguration
	if err := yaml.Unmarshal([]byte(content), &NewInitConfiguration); err != nil {
		logs.Lg.SysError("解析nacos配置", err, logs.Desc("解析nacos配置失败"))
	}
	return &NewInitConfiguration
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/21 3:39 下午
 * @Description: 返回配置文件自定义内容
**/
func ReadNacosSelfProfile(content string, pojo interface{}) {
	if err := yaml.Unmarshal([]byte(content), pojo); err != nil {
		logs.Lg.SysError("解析nacos配置", err, logs.Desc("解析nacos自定义配置失败"))
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/8 11:47 上午
 * @Description: 初始化限流配置
**/
func InitRateProfile() {
	if InitConfiguration.Rate.All { // 如果是全局，只设置一个值
		RateMap[constants.RATE_ALL], PastRateMap[constants.RATE_ALL] = ratelimit.New(InitConfiguration.Rate.Rate), InitConfiguration.Rate.Rate
	} else if len(InitConfiguration.Rate.InterfaceAndRate) != 0 { // 如果不是全局，那么逐个设置
		for k, v := range InitConfiguration.Rate.InterfaceAndRate {
			RateMap[addProgramName(k)], PastRateMap[addProgramName(k)] = ratelimit.New(v), v
		}
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/8 11:51 上午
 * @Description: 读取限流配置
**/
func ReadRateProfile(profile *InitNacosConfiguration) {
	if InitConfiguration.Rate.All { // 如果是全局，只设置一个值
		RateMap, PastRateMap = make(map[string]ratelimit.Limiter), make(map[string]int)
		RateMap[constants.RATE_ALL], PastRateMap[constants.RATE_ALL] = ratelimit.New(InitConfiguration.Rate.Rate), InitConfiguration.Rate.Rate
	} else if len(profile.Rate.InterfaceAndRate) != 0 { // 如果不是全局，那么逐个设置
		// 先删除取消的
		for k := range PastRateMap {
			if _, okay := profile.Rate.InterfaceAndRate[deleteProgramName(k)]; okay { // 存在
				continue
			} else { // 不存在
				delete(RateMap, k)
				delete(PastRateMap, k)
			}
		}
		// 再增加新增和修改的
		for k, v := range profile.Rate.InterfaceAndRate {
			if _, okay := PastRateMap[addProgramName(k)]; okay { // 存在
				if PastRateMap[addProgramName(k)] != v { // 如果不相同
					RateMap[addProgramName(k)], PastRateMap[addProgramName(k)] = ratelimit.New(v), v
				}
			} else { // 不存在
				RateMap[addProgramName(k)], PastRateMap[addProgramName(k)] = ratelimit.New(v), v
			}
		}
	} else { // 都没有，那么直接删除
		RateMap, PastRateMap = make(map[string]ratelimit.Limiter), make(map[string]int)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/8 2:12 下午
 * @Description: 在url中增加头路径
**/
func addProgramName(k string) string {
	k = fmt.Sprintf("/%s%s", InitConfiguration.Serve.ServerName, k)
	return k
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/13 10:15 上午
 * @Description: 在url中去除头路径
**/
func deleteProgramName(k string) string {
	k = k[strings.Index(k[1:], constants.BACKSLASH_MARK)+1:]
	return k
}
