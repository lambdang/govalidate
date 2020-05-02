package govalidate

import "errors"

const(
	REQUIRE = "require" //必须要有 查询时候用
	MAX = "max" //int最大值 字符串最大字符数
	MIN = "min" //int最小值
	CHOICE = "choice" //元素必须在合适的 1,2,3,4
	DATE = "date"
	LEN = "len" //str 长度
	PHONE = "phone" //手机号
	EMAIL = "email"  //email
)

type I interface {}
type M map[string]I
type Errors []error
var ER = errors.New