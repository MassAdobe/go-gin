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
	"github.com/MassAdobe/go-gin/config"
	"github.com/MassAdobe/go-gin/logs"
	"github.com/MassAdobe/go-gin/pojo"
	"github.com/MassAdobe/go-gin/systemUtils"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-26 21:08
 * @Description: 生成Access-Token
**/
func CreateToken(gUid int64, randomMark string) string {
	claim := jwt.MapClaims{
		config.TOKEN_VERIFY_KEY:   pojo.SysConf.TokenVerify,
		config.TOKEN_USER_KEY:     gUid,
		config.TOKEN_RANDOM_MARK:  randomMark,
		config.TOKEN_LOGIN_TM_KEY: systemUtils.RtnTmString(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	if tokenss, err := token.SignedString([]byte(pojo.SysConf.TokenVerify + pojo.SysConf.JwtKey)); err != nil {
		logs.Lg.Error("生成Access-Token错误", err)
	} else {
		split := strings.Split(tokenss, ".")
		tokenss = ""
		for idx, str := range split {
			if 0 != idx {
				split[idx] = str[10:] + str[:10]
				tokenss += "." + split[idx]
			} else {
				tokenss += split[idx]
			}
		}
		return tokenss
	}
	return ""
}
