/**
 * @Time : 2020-04-27 21:58
 * @Author : MassAdobe
 * @Description: pojo
**/
package pojo

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/15 2:02 下午
 * @Description: 用户基本信息
**/
type RequestUser struct {
	UserId   int64  `json:"user_id"`   // 用户ID
	UserFrom string `json:"user_from"` // 用户来源
}
