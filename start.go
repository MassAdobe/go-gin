/**
 * @Time : 2020/12/21 11:49 上午
 * @Author : MassAdobe
 * @Description: system
**/
package main

import (
	"context"
	"fmt"
	"github.com/MassAdobe/go-gin/db"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/nacos"
	"github.com/MassAdobe/go-gin/pojo"
	"github.com/MassAdobe/go-gin/systemUtils"
	"github.com/MassAdobe/go-gin/validated"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 1:49 下午
 * @Description: 预热项
**/
func init() {
	fmt.Println(fmt.Sprintf("【服务启动中...】"))
	s, _ := systemUtils.RunInLinuxWithErr("pwd")     // 执行linux命令获取当前路径
	sysData, _ := ioutil.ReadFile(s + "/config.yml") // 读取系统配置
	if err := yaml.Unmarshal(sysData, &pojo.InitConf); err != nil {
		fmt.Println(fmt.Sprintf("【MAIN】 %v %s", err, "解析系统配置失败"))
		os.Exit(1)
	}
	nacos.InitNacos()          // 初始化nacos配置
	nacos.NacosConfiguration() // nacos配置中心
	nacos.InitNacosProfile()   // 处理首次nacos获取到的配置信息
	nacos.NacosDiscovery()     // nacos服务注册发现
	logs.InitLogger(nacos.InitConfiguration.Log.Path,
		nacos.InitConfiguration.Serve.ServerName,
		nacos.InitConfiguration.Log.Level,
		nacos.InitConfiguration.Serve.Port) // 初始化日志
	db.InitDB()                             // 初始化DB
	validated.InitValidator()               // 初始化校验器
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 5:39 下午
 * @Description: 优雅停服
**/
func GracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGSTOP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, os.Interrupt)
	sig := <-quit
	logs.Lg.Info("准备关闭", logs.SpecDesc("收到信号", sig))
	now := time.Now()
	nacos.NacosDeregister() // nacos注销服务
	cxt, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := server.Shutdown(cxt); err != nil {
		logs.Lg.Error("关闭失败", err)
	}
	logs.Lg.Info("退出成功", logs.Desc(fmt.Sprintf("退出花费时间: %v", time.Since(now))))
}
