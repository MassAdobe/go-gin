/**
 * @Time : 2020-04-26 19:57
 * @Author : MassAdobe
 * @Description: config
**/
package constants

import "time"

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 2:02 下午
 * @Description: HTTP中的基本常量
**/
const (
	CONTENT_TYPE_KEY   = "Content-Type"                   // 请求协议种类键值
	CONTENT_TYPE_INNER = "application/json;charset=utf-8" // 请求协议种类内容
	REQUEST_USER_KEY   = "user"                           // 用户头信息键值
	REQUEST_TIMEOUT_TM = time.Duration(6 * time.Second)   // 五秒超时时间
	REQUEST_TRACE_ID   = "trace-id"                       // 跟踪ID
	REQUEST_STEP_ID    = "step-id"                        // 步骤ID
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/5 6:39 下午
 * @Description: 基本参数常量
**/
const (
	HAS_DELETED = "0" // 已删除
	NOT_DELETED = "1" // 未删除
	NOT_ENABLED = "0" // 无效
	IS_ENABLED  = "1" // 有效
)
