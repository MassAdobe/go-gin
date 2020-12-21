/**
 * @Time : 2020/12/17 4:24 下午
 * @Author : MassAdobe
 * @Description: nacos
**/
package nacos

import (
	"fmt"
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
		fmt.Println(fmt.Sprintf("【SYSTEM】%s %s %v %s", systemUtils.RtnCurTime(), "【nacos配置中心】", err, "读取nacos系统配置失败"))
		os.Exit(1)
	}
}
