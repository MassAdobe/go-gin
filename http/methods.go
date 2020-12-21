/**
 * @Time : 2020/12/18 3:00 下午
 * @Author : MassAdobe
 * @Description: http
**/
package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-gin/config"
	"go-gin/errs"
	"go-gin/logs"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 3:12 下午
 * @Description: Get请求
**/
func Get(ipPort, url string, params interface{}, c ...*gin.Context) ([]byte, error) {
	var (
		resp *http.Response
		err  error
	)
	url = url + urlEncode(params)
	client := http.Client{Timeout: config.REQUEST_TIMEOUT_TM}
	request, requestErr := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s%s", ipPort, url), nil)
	if nil != requestErr {
		logs.Lg.Error("Get请求", requestErr, c, logs.Desc(strings.NewReader(url)))
		return nil, requestErr
	}
	request.Header.Add(config.CONTENT_TYPE_KEY, config.CONTENT_TYPE_INNER)
	if len(c) != 0 {
		if user, has := c[0].Params.Get(config.REQUEST_USER_KEY); has {
			request.Header.Add(config.REQUEST_USER_KEY, user)
		}
		if trace, has := c[0].Params.Get(config.REQUEST_TRACE_ID); has {
			request.Header.Add(config.REQUEST_TRACE_ID, trace)
		}
		if step, has := c[0].Params.Get(config.REQUEST_STEP_ID); has {
			request.Header.Add(config.REQUEST_STEP_ID, step)
		}
	}
	if resp, err = client.Do(request); err != nil {
		defer resp.Body.Close()
		logs.Lg.Error("Get请求", err, c, logs.Desc(strings.NewReader(url)))
		return nil, err
	} else {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			logs.Lg.Error("Get请求", err, c, logs.Desc(resp.Body))
			return nil, err
		} else {
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
	var (
		resp *http.Response
		err  error
		req  string
	)
	if nil != params { // 判断参数是否为空
		if bytes, err := json.Marshal(params); err != nil {
			logs.Lg.Error("Post请求", err)
			panic(errs.NewError(errs.ErrJsonCode))
		} else {
			req = string(bytes)
		}
	}
	client := http.Client{Timeout: config.REQUEST_TIMEOUT_TM}
	requestUrl := fmt.Sprintf("http://%s%s", ipPort, url)
	request, requestErr := http.NewRequest(http.MethodPost, requestUrl, strings.NewReader(req))
	if nil != requestErr {
		logs.Lg.Error("Post请求", requestErr, c, logs.Desc(strings.NewReader(requestUrl)))
		return nil, requestErr
	}
	request.Header.Add(config.CONTENT_TYPE_KEY, config.CONTENT_TYPE_INNER)
	if len(c) != 0 {
		if user, has := c[0].Params.Get(config.REQUEST_USER_KEY); has {
			request.Header.Add(config.REQUEST_USER_KEY, user)
		}
		if trace, has := c[0].Params.Get(config.REQUEST_TRACE_ID); has {
			request.Header.Add(config.REQUEST_TRACE_ID, trace)
		}
		if step, has := c[0].Params.Get(config.REQUEST_STEP_ID); has {
			request.Header.Add(config.REQUEST_STEP_ID, step)
		}
	}
	if resp, err = client.Do(request); err != nil {
		defer resp.Body.Close()
		logs.Lg.Error("Post请求", err, c, logs.Desc(req))
		return nil, err
	} else {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			logs.Lg.Error("Post请求", err, c, logs.Desc(resp.Body))
			return nil, err
		} else {
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
	url = url + urlEncode(params)
	request, requestErr := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s%s", ipPort, url), nil)
	if nil != requestErr {
		logs.Lg.Error("Put请求", requestErr, c, logs.Desc(strings.NewReader(url)))
		return nil, requestErr
	}
	request.Header.Add(config.CONTENT_TYPE_KEY, config.CONTENT_TYPE_INNER)
	if len(c) != 0 {
		if user, has := c[0].Params.Get(config.REQUEST_USER_KEY); has {
			request.Header.Add(config.REQUEST_USER_KEY, user)
		}
		if trace, has := c[0].Params.Get(config.REQUEST_TRACE_ID); has {
			request.Header.Add(config.REQUEST_TRACE_ID, trace)
		}
		if step, has := c[0].Params.Get(config.REQUEST_STEP_ID); has {
			request.Header.Add(config.REQUEST_STEP_ID, step)
		}
	}
	client := http.Client{Timeout: config.REQUEST_TIMEOUT_TM}
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		logs.Lg.Error("Put请求", err, c, logs.Desc(strings.NewReader(url)))
	}
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		logs.Lg.Error("Put请求", err, c, logs.Desc(strings.NewReader(url)))
		return nil, err
	} else {
		return body, nil
	}
}

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/18 3:13 下午
 * @Description: Delete请求
**/
func Delete(ipPort, url string, params interface{}, c ...*gin.Context) ([]byte, error) {
	url = url + urlEncode(params)
	request, requestErr := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s%s", ipPort, url), nil)
	if nil != requestErr {
		logs.Lg.Error("Delete请求", requestErr, c, logs.Desc(strings.NewReader(url)))
		return nil, requestErr
	}
	request.Header.Add(config.CONTENT_TYPE_KEY, config.CONTENT_TYPE_INNER)
	if len(c) != 0 {
		if user, has := c[0].Params.Get(config.REQUEST_USER_KEY); has {
			request.Header.Add(config.REQUEST_USER_KEY, user)
		}
		if trace, has := c[0].Params.Get(config.REQUEST_TRACE_ID); has {
			request.Header.Add(config.REQUEST_TRACE_ID, trace)
		}
		if step, has := c[0].Params.Get(config.REQUEST_STEP_ID); has {
			request.Header.Add(config.REQUEST_STEP_ID, step)
		}
	}
	client := http.Client{Timeout: config.REQUEST_TIMEOUT_TM}
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		logs.Lg.Error("Delete请求", err, c, logs.Desc(strings.NewReader(url)))
	}
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		logs.Lg.Error("Delete请求", err, c, logs.Desc(strings.NewReader(url)))
		return nil, err
	} else {
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
	urlStr := strings.Join(urls, "&")
	return "?" + strings.ReplaceAll(urlStr, " ", ",")
}