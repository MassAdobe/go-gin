/**
 * @Time : 2020-04-30 11:36
 * @Author : MassAdobe
 * @Description: system
**/
package validated

import (
	"bytes"
	valid "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
	"go-gin/errs"
	"reflect"
)

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-30 15:45
 * @Description: 自定义验证器
**/
type Validator struct {
	Validate *validator.Validate
	Trans    ut.Translator
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-30 15:46
 * @Description: 全局验证器
**/
var GlobalValidator Validator

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-30 15:46
 * @Description: 初始化验证器
**/
func InitValidator() {
	zhs := valid.New()
	uni := ut.New(zhs, zhs)
	trans, ok := uni.GetTranslator("zh")
	if !ok {
		panic(errs.NewMsgError(errs.ErrSystemCode, "初始化验证器错误"))
	}
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string { // 收集结构体中的comment标签，用于替换英文字段名称，这样返回错误就能展示中文字段名称了
		return fld.Tag.Get("comment")
	})
	err := zh.RegisterDefaultTranslations(validate, trans) // 注册中文翻译
	if err != nil {
		panic(errs.NewMsgError(errs.ErrSystemCode, err.Error()))
	}
	GlobalValidator.Validate = validate
	GlobalValidator.Trans = trans
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-30 15:47
 * @Description: 验证器通用验证方法
**/
func (m *Validator) Check(value interface{}) error {
	err := m.Validate.Struct(value) // 首先使用validator.v10进行验证
	if err != nil {
		verrs, ok := err.(validator.ValidationErrors)
		if !ok { // 几乎不会出现，除非验证器本身异常无法转换，以防万一就判断一下好了
			panic(errs.NewError(errs.ErrParamCheckingCode))
		}
		errBuf := bytes.Buffer{} // 将所有的参数错误进行翻译然后拼装成字符串返回
		for i := 0; i < len(verrs); i++ {
			errBuf.WriteString(verrs[i].Translate(m.Trans) + "\n")
		}
		errStr := errBuf.String() // 删除掉最后一个空格和换行符
		panic(errs.NewMsgError(errs.ErrParamCode, errStr[:len(errStr)-1]))
	}
	if v, ok := value.(CanCheck); ok { // 如果它实现了CanCheck接口，就进行自定义验证
		return v.Check()
	}
	return nil
}

/**
 * @Author: MassAdobe
 * @TIME: 2020-04-30 15:51
 * @Description: 如果需要特殊校验，可以实现验证接口，或者通过自定义tag标签实现
**/
type CanCheck interface {
	Check() error
}
