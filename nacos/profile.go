/**
 * @Time : 2020/12/17 4:24 下午
 * @Author : MassAdobe
 * @Description: nacos
**/
package nacos

import (
	"fmt"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/systemUtils"
	"gopkg.in/yaml.v2"
	"os"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 4:24 下午
 * @Description: 处理首次nacos获取到的配置信息
**/
func InitNacosProfile() {
	if err := yaml.Unmarshal([]byte(NacosContent), &InitConfiguration); err != nil {
		fmt.Println(fmt.Sprintf(`{"log_level":"INFO","time":"%s","msg":"%s","server_name":"%s","desc":"%s"}`, systemUtils.RtnCurTime(), "配置中心", "未知", "读取nacos系统配置失败"))
		os.Exit(1)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/21 3:05 下午
 * @Description: 返回配置文件内容
**/
func ReadNacosProfile(content string) *InitNacosConfiguration {
	if err := yaml.Unmarshal([]byte(content), &InitConfiguration); err != nil {
		logs.Lg.SysError("解析nacos配置", err, logs.Desc("解析nacos配置失败"))
	}
	return &InitConfiguration
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
