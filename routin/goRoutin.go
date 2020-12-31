/**
 * @Time : 2020/12/31 3:04 下午
 * @Author : MassAdobe
 * @Description: routin
**/
package routin

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020/12/31 2:55 下午
 * @Description: 获取goroutine的ID
**/
func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
