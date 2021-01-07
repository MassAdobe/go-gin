/**
 * @Time : 2020/12/18 2:29 下午
 * @Author : MassAdobe
 * @Description: http
**/
package http

import (
	"fmt"
	"github.com/MassAdobe/go-gin/errs"
	"github.com/MassAdobe/go-gin/goContext"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/nacos"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 2:30 下午
 * @Description: 服务内部调用get请求结构体
**/
type FeignRequest struct {
	Body       interface{}        // 请求参数，可以为空
	ServerName string             // 服务名，不能为空
	GroupName  string             // 组别名，不能为空
	Url        string             // 调用URL(二级路径)
	C          *goContext.Context // 当前请求的上下文
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 2:30 下午
 * @Description: 服务内部调用get请求
**/
func (this *FeignRequest) FeignGet() (feign []byte, err error) {
	if 0 == nacos.InitConfiguration.Feign.RetryNum {
		if feign, err = getFeign(this.ServerName, this.GroupName, this.Url, this.Body, this.C); err != nil {
			panic(errs.NewError(errs.ErrGetRequestCode, err))
		} else {
			return
		}
	} else {
		for i := 0; i < nacos.InitConfiguration.Feign.RetryNum; i++ {
			if feign, err = getFeign(this.ServerName, this.GroupName, this.Url, this.Body, this.C); err == nil {
				return
			}
		}
	}
	this.C.SysError("服务内部调用get请求", err, logs.Desc(fmt.Sprintf("超过设置调用次数: %d", nacos.InitConfiguration.Feign.RetryNum)))
	return
}

func getFeign(serviceName, groupName, url string, params interface{}, c *goContext.Context) ([]byte, error) {
	if instance, err := nacos.NacosGetServer(serviceName, groupName); err != nil {
		c.SysError("服务内部调用get请求", err)
		panic(errs.NewError(errs.ErrInnerCallingCode, err))
	} else {
		if rtn, err := Get(fmt.Sprintf("%s:%d", instance.Ip, instance.Port), url, params, c.GinContext); err != nil {
			c.SysError("服务内部调用get请求", err)
			return nil, err
		} else {
			return rtn, nil
		}
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 2:30 下午
 * @Description: 服务内部调用post请求
**/
func (this *FeignRequest) FeignPost() (feign []byte, err error) {
	if 0 == nacos.InitConfiguration.Feign.RetryNum {
		if feign, err = postFeign(this.ServerName, this.GroupName, this.Url, this.Body, this.C); err != nil {
			panic(errs.NewError(errs.ErrGetRequestCode, err))
		} else {
			return
		}
	} else {
		for i := 0; i < nacos.InitConfiguration.Feign.RetryNum; i++ {
			if feign, err = postFeign(this.ServerName, this.GroupName, this.Url, this.Body, this.C); err == nil {
				return
			}
		}
	}
	this.C.SysError("服务内部调用post请求", err, logs.Desc(fmt.Sprintf("超过设置调用次数: %d", nacos.InitConfiguration.Feign.RetryNum)))
	return
}

func postFeign(serviceName, groupName, url string, params interface{}, c *goContext.Context) ([]byte, error) {
	if instance, err := nacos.NacosGetServer(serviceName, groupName); err != nil {
		c.SysError("服务内部调用post请求", err)
		panic(errs.NewError(errs.ErrInnerCallingCode, err))
	} else {
		if rtn, err := Post(fmt.Sprintf("%s:%d", instance.Ip, instance.Port), url, params, c.GinContext); err != nil {
			c.SysError("服务内部调用post请求", err)
			return nil, err
		} else {
			return rtn, nil
		}
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 5:21 下午
 * @Description: 服务内部调用put请求
**/
func (this *FeignRequest) FeignPut() (feign []byte, err error) {
	if 0 == nacos.InitConfiguration.Feign.RetryNum {
		if feign, err = putFeign(this.ServerName, this.GroupName, this.Url, this.Body, this.C); err != nil {
			panic(errs.NewError(errs.ErrGetRequestCode, err))
		} else {
			return
		}
	} else {
		for i := 0; i < nacos.InitConfiguration.Feign.RetryNum; i++ {
			if feign, err = putFeign(this.ServerName, this.GroupName, this.Url, this.Body, this.C); err == nil {
				return
			}
		}
	}
	this.C.SysError("服务内部调用put请求", err, logs.Desc(fmt.Sprintf("超过设置调用次数: %d", nacos.InitConfiguration.Feign.RetryNum)))
	return
}

func putFeign(serviceName, groupName, url string, params interface{}, c *goContext.Context) ([]byte, error) {
	if instance, err := nacos.NacosGetServer(serviceName, groupName); err != nil {
		c.SysError("服务内部调用put请求", err)
		panic(errs.NewError(errs.ErrInnerCallingCode, err))
	} else {
		if rtn, err := Put(fmt.Sprintf("%s:%d", instance.Ip, instance.Port), url, params, c.GinContext); err != nil {
			c.SysError("服务内部调用put请求", err)
			return nil, err
		} else {
			return rtn, nil
		}
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 5:21 下午
 * @Description: 服务内部调用Delete请求
**/
func (this *FeignRequest) FeignDelete() (feign []byte, err error) {
	if 0 == nacos.InitConfiguration.Feign.RetryNum {
		if feign, err = deleteFeign(this.ServerName, this.GroupName, this.Url, this.Body, this.C); err != nil {
			panic(errs.NewError(errs.ErrGetRequestCode, err))
		} else {
			return
		}
	} else {
		for i := 0; i < nacos.InitConfiguration.Feign.RetryNum; i++ {
			if feign, err = deleteFeign(this.ServerName, this.GroupName, this.Url, this.Body, this.C); err == nil {
				return
			}
		}
	}
	this.C.SysError("服务内部调用delete请求", err, logs.Desc(fmt.Sprintf("超过设置调用次数: %d", nacos.InitConfiguration.Feign.RetryNum)))
	return
}

func deleteFeign(serviceName, groupName, url string, params interface{}, c *goContext.Context) ([]byte, error) {
	if instance, err := nacos.NacosGetServer(serviceName, groupName); err != nil {
		c.SysError("服务内部调用delete请求", err)
		panic(errs.NewError(errs.ErrInnerCallingCode, err))
	} else {
		if rtn, err := Delete(fmt.Sprintf("%s:%d", instance.Ip, instance.Port), url, params, c.GinContext); err != nil {
			c.SysError("服务内部调用delete请求", err)
			return nil, err
		} else {
			return rtn, nil
		}
	}
}
