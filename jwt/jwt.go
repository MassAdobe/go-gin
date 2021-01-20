/**
 * @Time : 2020-04-26 19:57
 * @Author : MassAdobe
 * @Description: system
**/
package jwt

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:08
 * @Description: JWT工具类
**/
import (
	"github.com/MassAdobe/go-gin/constants"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/nacos"
	"github.com/MassAdobe/go-gin/systemUtils"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/15 1:54 下午
 * @Description: 生成Access-Token
**/
func CreateToken(userId int64, userFrom string) string {
	claim := jwt.MapClaims{
		constants.TOKEN_USER_KEY:     userId,
		constants.TOKEN_USER_FROM:    userFrom,
		constants.TOKEN_LOGIN_TM_KEY: systemUtils.RtnTmString(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	if tokenss, err := token.SignedString([]byte(nacos.InitConfiguration.AccessToken.Verify)); err != nil {
		logs.Lg.Error("生成Access-Token错误", err)
	} else {
		split := strings.Split(tokenss, constants.FULL_STOP_MARK)
		tokenss = ""
		for idx, str := range split {
			if 0 != idx {
				split[idx] = str[10:] + str[:10]
				tokenss += constants.FULL_STOP_MARK + split[idx]
			} else {
				tokenss += split[idx]
			}
		}
		return tokenss
	}
	return ""
}
