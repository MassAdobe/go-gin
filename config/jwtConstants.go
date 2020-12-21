/**
 * @Time : 2020-04-26 19:57
 * @Author : MassAdobe
 * @Description: config
**/
package config

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:06
 * @Description: jwt常量
**/
const (
	TOKEN_VERIFY_KEY   = "v_scrt" // Token中校验元素secret
	TOKEN_USER_KEY     = "g_uid"  // Token中的用户KEY
	TOKEN_LOGIN_TM_KEY = "lgn_tm" // Token中的Login时间
	TOKEN_RANDOM_MARK  = "ran_mk" // 前端的随机校验码
)
