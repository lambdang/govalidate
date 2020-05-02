package govalidate

import (
	"testing"
	"fmt"
)


type A struct{
	Name string `json:"name" valid:"require;max=2;min=2" name:"姓名"`
	Age int`json:"age" valid:"require;max=11;min=2" name:"年龄"`
	Level string `json:"level" valid:"choice=1,2,3"`
}


func TestName(t *testing.T) {

	a:=A{}
	a.Name = "d"
	a.Age = 1
	a.Level = "8"
	errors,ok:=ValidStruct(&a)
	//返回 ok 校验是否成功 errors是map[string]string 字典key属性名或者err
	//tag 中name为验证结果所用的名字，如果不定义则用属性名
	fmt.Println(errors,ok)
	//map[age:年龄不能小于2 level:Level必须是1,2,3其中一个 name:姓名长度不能小于2] false

	//默认注册的验证有
	//REQUIRE = "require" //必须要有 查询时候用
	//MAX = "max" //int最大值 字符串最大字符数
	//MIN = "min" //int最小值
	//CHOICE = "choice" //元素必须在合适的 1,2,3,4
	//DATE = "date"
	//LEN = "len" //str 长度
	//PHONE = "phone" //手机号
	//EMAIL = "email"  //email
	//需要其他验证可自行添加
}
