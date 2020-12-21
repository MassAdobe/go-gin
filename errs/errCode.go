/**
 * @Time : 2020-04-26 19:57
 * @Author : MassAdobe
 * @Description: config
**/
package errs

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:05
 * @Description: 错误封装
**/
const (
	/*error code*/
	SuccessCode                        = 0    // 成功
	ErrSystemCode                      = iota // 内部错误
	ErrNotFoundCode                           // 资源没有找到
	ErrParamCode                              // 输入参数错误
	ErrJsonCode                               // json解析错误
	ErrRedisCode                              // redis错误
	ErrHeaderUserCode                         // 用户基本信息错误
	ErrParamCheckingCode                      // 参数校验错误
	ErrDataBaseCode                           // 数据库错误
	ErrInnerCallingCode                       // 内部调用获取服务错误
	ErrInnerCallingServiceNotFoundCode        // 内部调用服务目标服务不存在错误
	ErrInnerCallingExecCode                   // 内部调用服务错误
	ErrPostRequestCode                        // Post请求错误
	ErrGetRequestCode                         // Get请求错误
	ErrPutRequestCode                         // Put请求错误
	ErrDeleteRequestCode                      // Deleted请求错误

	//自定义错误码
	ErrLoginCode = 10000 + iota // 登录错误

	/*error desc*/
	SuccessDesc                        = "成功"
	ErrSystemDesc                      = "内部错误"
	ErrNotFoundDesc                    = "资源没有找到"
	ErrParamDesc                       = "输入参数错误"
	ErrJsonDesc                        = "Json转换错误"
	ErrRedisDesc                       = "redis错误"
	ErrHeaderUserDesc                  = "用户基本信息错误"
	ErrParamCheckingDesc               = "参数校验错误"
	ErrDataBaseDesc                    = "数据库错误"
	ErrInnerCallingDesc                = "内部调用获取服务错误"
	ErrInnerCallingServiceNotFoundDesc = "内部调用服务目标服务不存在错误"
	ErrInnerCallingExecDesc            = "内部调用服务错误"
	ErrPostRequestDesc                 = "Post请求错误"
	ErrGetRequestDesc                  = "Get请求错误"
	ErrPutRequestDesc                  = "Put请求错误"
	ErrDeleteRequestDesc               = "Deleted请求错误"

	//自定义错误描述
	ErrLoginDesc = "登录错误(用户名密码错误或不存在相关用户)"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:06
 * @Description: 错误参数体
**/
var CodeDescMap = map[int]string{
	// 系统错误
	SuccessCode:                        SuccessDesc,
	ErrSystemCode:                      ErrSystemDesc,
	ErrNotFoundCode:                    ErrNotFoundDesc,
	ErrParamCode:                       ErrParamDesc,
	ErrJsonCode:                        ErrJsonDesc,
	ErrRedisCode:                       ErrRedisDesc,
	ErrHeaderUserCode:                  ErrHeaderUserDesc,
	ErrParamCheckingCode:               ErrParamCheckingDesc,
	ErrDataBaseCode:                    ErrDataBaseDesc,
	ErrInnerCallingCode:                ErrInnerCallingDesc,
	ErrInnerCallingServiceNotFoundCode: ErrInnerCallingServiceNotFoundDesc,
	ErrInnerCallingExecCode:            ErrInnerCallingExecDesc,
	ErrPostRequestCode:                 ErrPostRequestDesc,
	ErrGetRequestCode:                  ErrGetRequestDesc,
	ErrPutRequestCode:                  ErrPutRequestDesc,
	ErrDeleteRequestCode:               ErrDeleteRequestDesc,
	// 自定义错误
	ErrLoginCode: ErrLoginDesc,
}
