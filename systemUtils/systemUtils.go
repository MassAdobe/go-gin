/**
 * @Time : 2020-04-26 17:38
 * @Author : MassAdobe
 * @Description: 基本服务工具类
**/
package systemUtils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go-gin/errs"
	"go-gin/logs"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 3:35 下午
 * @Description: 常量池
**/
const (
	TimeFormatMS    = "2006-01-02 15:04:05"
	TimeFormatMonth = "2006-01-02"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 6:55 下午
 * @Description: 生成随机时间戳标志位
**/
func RandomTimestampMark() string {
	return fmt.Sprintf("%d%d",
		time.Now().UnixNano(),
		RandInt64(1000, 9999))
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 3:36 下午
 * @Description: 运行当前系统命令
**/
func RunInLinuxWithErr(cmd string) (string, error) {
	result, err := exec.Command(cmd).Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	return strings.TrimSpace(string(result)), err
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:12
 * @Description: md5加密
**/
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:13
 * @Description: 区间随机数；返回int64
**/
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max-min) + min
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:13
 * @Description: 区间随机数；返回int
**/
func RandInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:13
 * @Description: 随机字符串基础值
**/
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:13
 * @Description: 随机字符串
**/
func RandSeq(n int) string {
	b, r := make([]rune, n), rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-05-29 21:47
 * @Description: 生成手机验证码
**/
func RandCodeSeq() string {
	var codes = []rune("0123456789")
	b, r := make([]rune, 6), rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = codes[r.Intn(len(codes))]
	}
	return string(b)
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:13
 * @Description: 获取当前IP地址
**/
func GetIntranetIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-27 21:19
 * @Description: 转化结构体为JSON字符串
**/
func Marshal(pojo interface{}) string {
	if bytes, err := json.Marshal(pojo); err != nil {
		logs.Lg.Error("结构体转JSON格式错误", err)
		panic(errs.NewError(errs.ErrJsonCode))
	} else {
		return string(bytes)
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 4:33 下午
 * @Description: 获取当前系统IP
**/
func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("没链接网络")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}
	return ip
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/21 10:53 上午
 * @Description: 获取请求url上的所有参数
**/
func GetRequestUrlParams(uri string) string {
	if strings.Contains(uri, "?") {
		return uri[strings.Index(uri, "?")+1:]
	}
	return ""
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:13
 * @Description: 返回时间字符串
**/
func RtnTmString() (timsStr string) {
	timsStr = time.Now().Format(TimeFormatMS)
	return
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 4:54 下午
 * @Description: 返回当前时间戳
**/
func RtnCurTime() string {
	return time.Now().Format(TimeFormatMS)
}