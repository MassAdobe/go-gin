/**
 * @Time : 2021/1/15 3:06 下午
 * @Author : MassAdobe
 * @Description: go_gin
**/
package main

import (
	"fmt"
	"strings"
)

func main() {
	token := "12345678987654321.12345678987654321.12345678987654321.12345678987654321"
	token1 := encodeToken(token)
	fmt.Println(token1)
	token2 := decodeToken(token1)
	fmt.Println(token2)
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/15 2:58 下午
 * @Description: 加密token
**/
func encodeToken(token string) string {
	split, rtn := strings.Split(token, "."), ""
	for idx, str := range split {
		if 0 != idx {
			split[idx] = str[10:] + str[:10]
			rtn += "." + split[idx]
		} else {
			rtn += split[idx]
		}
	}
	return rtn
}

/**
 * @Author: MassAdobe
 * @TIME: 2021/1/15 2:58 下午
 * @Description: 解密token
**/
func decodeToken(token string) string {
	split, rtn := strings.Split(token, "."), ""
	for idx, str := range split {
		if 0 != idx {
			split[idx] = str[len(str)-10:] + str[:len(str)-10]
			rtn += "." + split[idx]
		} else {
			rtn += split[idx]
		}
	}
	return rtn
}
