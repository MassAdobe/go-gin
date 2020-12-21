/**
 * @Time : 2020-04-26 19:57
 * @Author : MassAdobe
 * @Description: config
**/
package config

import "time"

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/17 2:02 下午
 * @Description: HTTP中的基本常量
**/
const (
	CONTENT_TYPE_KEY   = "Content-Type"                 // 请求协议种类键值
	CONTENT_TYPE_INNER = "application/json"             // 请求协议种类内容
	REQUEST_USER_KEY   = "user"                         // 用户头信息键值
	REQUEST_TIMEOUT_TM = time.Duration(5 * time.Second) // 五秒超时时间
	REQUEST_TRACE_ID   = "trace-id"                     // 跟踪ID
	REQUEST_STEP_ID    = "step-id"                      // 步骤ID
)

