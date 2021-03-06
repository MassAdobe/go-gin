/**
 * @Time : 2020/12/18 3:00 下午
 * @Author : MassAdobe
 * @Description: http
**/
package http

import (
	"encoding/json"
	"fmt"
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/errs"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"
)

var (
	client *http.Client
	req    string
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/15 10:46 上午
 * @Description: 初始化连接池
**/
func init() {
	client = &http.Client{
		Timeout: constants.REQUEST_TIMEOUT_TM,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // 连接超时
				KeepAlive: 30 * time.Second, // 长连接超时时间
			}).DialContext,
			MaxIdleConns:          50,               // 最大空闲连接
			IdleConnTimeout:       90 * time.Second, // 空闲超时时间
			TLSHandshakeTimeout:   10 * time.Second, // tls握手超时时间
			ExpectContinueTimeout: 1 * time.Second,  // 100-continue 超时时间
		},
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/15 10:50 上午
 * @Description: 关闭HTTP连接池
**/
func CloseHttpConnectionPool() {
	client.CloseIdleConnections()
	logs.Lg.SysDebug("http连接池", logs.Desc("关闭http最大空闲连接"))
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 3:12 下午
 * @Description: Get请求
**/
func Get(ipPort, url string, params interface{}, c ...*gin.Context) ([]byte, error) {
	logs.Lg.SysDebug("Get请求", logs.Desc("进入Get请求方法"))
	url = url + urlEncode(params)
	request, requestErr := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s%s", ipPort, url), nil)
	if nil != requestErr {
		logs.Lg.SysError("Get请求", requestErr, c, logs.Desc(strings.NewReader(url)))
		return nil, requestErr
	}
	logs.Lg.SysDebug("Get请求", c, logs.Desc(fmt.Sprintf("请求地址：%s", request.RequestURI)))
	request.Header.Add(constants.CONTENT_TYPE_KEY, constants.CONTENT_TYPE_INNER)
	if len(c) != 0 {
		if user, has := c[0].Params.Get(constants.REQUEST_USER_KEY); has {
			request.Header.Add(constants.REQUEST_USER_KEY, user)
		}
		if trace, has := c[0].Params.Get(constants.REQUEST_TRACE_ID); has {
			request.Header.Add(constants.REQUEST_TRACE_ID, trace)
		}
		if step, has := c[0].Params.Get(constants.REQUEST_STEP_ID); has {
			request.Header.Add(constants.REQUEST_STEP_ID, step)
		}
	}
	var (
		resp *http.Response
		err  error
	)
	defer func() {
		if respErr := resp.Body.Close(); respErr != nil {
			logs.Lg.SysError("Get请求", respErr, logs.Desc("关闭返回体失败"))
			return
		}
		logs.Lg.SysDebug("Get请求", logs.Desc("关闭返回体成功"))
	}()
	if resp, err = client.Do(request); err != nil {
		logs.Lg.SysError("Get请求", err, c, logs.Desc(strings.NewReader(url)))
		return nil, err
	} else {
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			logs.Lg.SysError("Get请求", err, c, logs.Desc(resp.Body))
			return nil, err
		} else {
			logs.Lg.SysDebug("Get请求", c, logs.Desc(fmt.Sprintf("请求成功，返回：%s", string(body))))
			return body, nil
		}
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 3:09 下午
 * @Description: Post请求
**/
func Post(ipPort, url, params interface{}, c ...*gin.Context) ([]byte, error) {
	logs.Lg.SysDebug("Post请求", logs.Desc("进入Post请求方法"))
	if nil != params { // 判断参数是否为空
		if bytes, err := json.Marshal(params); err != nil {
			logs.Lg.SysError("Post请求", err)
			panic(errs.NewError(errs.ErrJsonCode))
		} else {
			req = string(bytes)
		}
	}
	requestUrl := fmt.Sprintf("http://%s%s", ipPort, url)
	request, requestErr := http.NewRequest(http.MethodPost, requestUrl, strings.NewReader(req))
	if nil != requestErr {
		logs.Lg.SysError("Post请求", requestErr, c, logs.Desc(strings.NewReader(requestUrl)))
		return nil, requestErr
	}
	logs.Lg.SysDebug("Post请求", c, logs.Desc(fmt.Sprintf("请求地址：%s", request.RequestURI)))
	request.Header.Add(constants.CONTENT_TYPE_KEY, constants.CONTENT_TYPE_INNER)
	if len(c) != 0 {
		if user, has := c[0].Params.Get(constants.REQUEST_USER_KEY); has {
			request.Header.Add(constants.REQUEST_USER_KEY, user)
		}
		if trace, has := c[0].Params.Get(constants.REQUEST_TRACE_ID); has {
			request.Header.Add(constants.REQUEST_TRACE_ID, trace)
		}
		if step, has := c[0].Params.Get(constants.REQUEST_STEP_ID); has {
			request.Header.Add(constants.REQUEST_STEP_ID, step)
		}
	}
	var (
		resp *http.Response
		err  error
	)
	defer func() {
		if respErr := resp.Body.Close(); respErr != nil {
			logs.Lg.SysError("Post请求", respErr, logs.Desc("关闭返回体失败"))
			return
		}
		logs.Lg.SysDebug("Post请求", logs.Desc("关闭返回体成功"))
	}()
	if resp, err = client.Do(request); err != nil {
		logs.Lg.SysError("Post请求", err, c, logs.Desc(req))
		return nil, err
	} else {
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			logs.Lg.SysError("Post请求", err, c, logs.Desc(resp.Body))
			return nil, err
		} else {
			logs.Lg.SysDebug("Post请求", c, logs.Desc(fmt.Sprintf("请求成功，返回：%s", string(body))))
			return body, nil
		}
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 3:13 下午
 * @Description: Put请求
**/
func Put(ipPort, url string, params interface{}, c ...*gin.Context) ([]byte, error) {
	logs.Lg.SysDebug("Put请求", logs.Desc("进入Put请求方法"))
	url = url + urlEncode(params)
	request, requestErr := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s%s", ipPort, url), nil)
	if nil != requestErr {
		logs.Lg.SysError("Put请求", requestErr, c, logs.Desc(strings.NewReader(url)))
		return nil, requestErr
	}
	logs.Lg.SysDebug("Put请求", c, logs.Desc(fmt.Sprintf("请求地址：%s", request.RequestURI)))
	request.Header.Add(constants.CONTENT_TYPE_KEY, constants.CONTENT_TYPE_INNER)
	if len(c) != 0 {
		if user, has := c[0].Params.Get(constants.REQUEST_USER_KEY); has {
			request.Header.Add(constants.REQUEST_USER_KEY, user)
		}
		if trace, has := c[0].Params.Get(constants.REQUEST_TRACE_ID); has {
			request.Header.Add(constants.REQUEST_TRACE_ID, trace)
		}
		if step, has := c[0].Params.Get(constants.REQUEST_STEP_ID); has {
			request.Header.Add(constants.REQUEST_STEP_ID, step)
		}
	}
	resp, err := client.Do(request)
	defer func() {
		if respErr := resp.Body.Close(); respErr != nil {
			logs.Lg.SysError("Put请求", respErr, logs.Desc("关闭返回体失败"))
			return
		}
		logs.Lg.SysDebug("Put请求", logs.Desc("关闭返回体成功"))
	}()
	if err != nil {
		logs.Lg.SysError("Put请求", err, c, logs.Desc(strings.NewReader(url)))
	}
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		logs.Lg.SysError("Put请求", err, c, logs.Desc(strings.NewReader(url)))
		return nil, err
	} else {
		logs.Lg.SysDebug("Put请求", c, logs.Desc(fmt.Sprintf("请求成功，返回：%s", string(body))))
		return body, nil
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 3:13 下午
 * @Description: Delete请求
**/
func Delete(ipPort, url string, params interface{}, c ...*gin.Context) ([]byte, error) {
	logs.Lg.SysDebug("Delete请求", logs.Desc("进入Post请求方法"))
	url = url + urlEncode(params)
	request, requestErr := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s%s", ipPort, url), nil)
	if nil != requestErr {
		logs.Lg.SysError("Delete请求", requestErr, c, logs.Desc(strings.NewReader(url)))
		return nil, requestErr
	}
	logs.Lg.SysDebug("Delete请求", c, logs.Desc(fmt.Sprintf("请求地址：%s", request.RequestURI)))
	request.Header.Add(constants.CONTENT_TYPE_KEY, constants.CONTENT_TYPE_INNER)
	if len(c) != 0 {
		if user, has := c[0].Params.Get(constants.REQUEST_USER_KEY); has {
			request.Header.Add(constants.REQUEST_USER_KEY, user)
		}
		if trace, has := c[0].Params.Get(constants.REQUEST_TRACE_ID); has {
			request.Header.Add(constants.REQUEST_TRACE_ID, trace)
		}
		if step, has := c[0].Params.Get(constants.REQUEST_STEP_ID); has {
			request.Header.Add(constants.REQUEST_STEP_ID, step)
		}
	}
	resp, err := client.Do(request)
	defer func() {
		if respErr := resp.Body.Close(); respErr != nil {
			logs.Lg.SysError("Delete请求", respErr, logs.Desc("关闭返回体失败"))
			return
		}
		logs.Lg.SysDebug("Delete请求", logs.Desc("关闭返回体成功"))
	}()
	if err != nil {
		logs.Lg.SysError("Delete请求", err, c, logs.Desc(strings.NewReader(url)))
	}
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		logs.Lg.SysError("Delete请求", err, c, logs.Desc(strings.NewReader(url)))
		return nil, err
	} else {
		logs.Lg.SysDebug("Delete请求", c, logs.Desc(fmt.Sprintf("请求成功，返回：%s", string(body))))
		return body, err
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 4:31 下午
 * @Description: get,put,delete请求转
**/
func urlEncode(params interface{}) string {
	if nil == params {
		return ""
	}
	urls := make([]string, 0)
	t := reflect.TypeOf(params)
	v := reflect.ValueOf(params)
	for k := 0; k < t.NumField(); k++ {
		urls = append(urls, fmt.Sprintf("%s=%v", t.Field(k).Tag.Get("json"), v.Field(k).Interface()))
	}
	urlStr := strings.Join(urls, constants.AND_MARK)
	return constants.QUESTION_MARK + strings.ReplaceAll(urlStr, constants.SPACE_MARK, constants.COMMA_MARK)
}
