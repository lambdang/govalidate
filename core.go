package govalidate

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/**
规则解析器
@params rulestr 规则字符串，一条规则 通过tag获取 ""
@return map[string]string 规则集
*/
func ParseRule(ruleStr string) (map[string]string, error) {
	var ruleMap = map[string]string{}
	ruleL := strings.Split(ruleStr, ";")
	for _, v := range ruleL {
		_t := strings.Split(v, "=")
		_l := len(_t)
		if _l == 0 {
			continue
		} else if _l == 1 {
			ruleMap[_t[0]] = "true"
		} else if _l == 2 {
			ruleMap[_t[0]] = _t[1]
		} else {
			return nil, fmt.Errorf("规则配置错误!")
		}
	}

	return ruleMap, nil
}

//=======================================================
/**
规则处理器
 */
func NewValidStore() ValidStore {
	v := ValidStore{
		//注册验证用基本函数
		CHOICE: validChoice,
		MAX:    validMax,
		MIN:    validMin,
		DATE:   validDate,
		LEN:    validLen,
		PHONE:  validPhone,
		EMAIL:  validEmail,
	}
	return v
}

/*
reflect.value
 */
type validFun func(interface{}, string, string) error

type ValidStore map[string]validFun

// 添加规则
func (this ValidStore) AddValidFun(name string, f validFun) {
	this[name] = f
}

/*
验证 单条数据

@params raw interface{} 验证的数据

@params name string 验证数据的名字

@params rule map[string]string 待验证的数据 key为需要使用的规则名 value为要验证的数据

return []error 未过检信息

reutrn bool 是否通过 true 通过 false 未通过

验证单条 有错误则返回
 */
func (this ValidStore) ValidRule(rv reflect.Value, name string, rule map[string]string) (error, bool) {
	var (
		//rerr errorerror
		ok  bool
		err error
		va  string
	)
	_, ok = rule[REQUIRE]
	raw := rv.Interface()

	//'必须有'验证 不能为空或者为""
	if raw == nil {
		if ok {
			return fmt.Errorf("缺少%s", name), false
		}
	} else if v, _ok := raw.(string); _ok {
		if v == "" {
			if ok {
				return fmt.Errorf("缺少%s", name), false
			}else{
				return nil,true
			}
		}
	} else if v, _ok := raw.(float64); _ok {
		if v == 0 {
			if ok {
				return fmt.Errorf("缺少%s", name), false
			}else{
				return nil,true
			}
		}
	} else if v, _ok := raw.(int); _ok {
		if v == 0 {
			if ok {
				return fmt.Errorf("缺少%s", name), false
			}else{
				return nil,true
			}
		}
	} else if v, _ok := raw.(int64); _ok {
		if v == 0 {
			if ok {
				return fmt.Errorf("缺少%s", name), false
			}else{
				return nil,true
			}
		}
	}

	for k, f := range this {
		if va, ok = rule[k]; ok {
			if err = f(raw, va, name); err != nil {
				return err, false
			}
		}
	}
	return nil, true
}

////======================判断规则=========================


//只支持数字和字符串
func validChoice(value interface{}, rule string, name string) error {
	d := strings.Split(rule, ",")
	var valueStr string
	switch value.(type) {
	case int, int8, int16, int32, int64:
		valueStr = strconv.FormatInt(reflect.ValueOf(1).Int(), 10)
		break
	case uint, uint8, uint16, uint32, uint64:
		valueStr = strconv.FormatUint(reflect.ValueOf(value).Uint(), 10)
		break
	case float32, float64:
		valueStr = strconv.FormatFloat(reflect.ValueOf(value).Float(), '6', -1, 64)
		break
	case string:
		valueStr = value.(string)
	}
	for _, v := range d {
		if valueStr == v {
			return nil
		}
	}
	return fmt.Errorf("%s必须是%s其中一个", name, rule)

}

//最大 支持数字和字符长度
func validMax(value interface{}, rule, name string) error {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		_tmp := reflect.ValueOf(value).Int()
		if _maxNum, err := strconv.ParseInt(rule, 10, 64); err != nil {
			return fmt.Errorf("%s配置错误", rule)
		} else {
			if (_tmp > _maxNum) {
				return fmt.Errorf("%s不能大于%s", name, rule)
			}
		}
		break
	case uint, uint8, uint16, uint32, uint64:
		_tmp := reflect.ValueOf(value).Uint()
		if _maxNum, err := strconv.ParseUint(rule, 10, 64); err != nil {
			return fmt.Errorf("%s配置错误", rule)
		} else {
			if (_tmp > _maxNum) {
				return fmt.Errorf("%s不能大于%s", name, rule)
			}
		}

	case float32, float64:
		_tmp := reflect.ValueOf(value).Float()
		if _maxNum, err := strconv.ParseFloat(rule, 64); err != nil {
			return fmt.Errorf("%s配置错误", rule)
		} else {
			if (_tmp > _maxNum) {
				return fmt.Errorf("%s不能大于%s", name, rule)
			}
		}
	case string:
		_tmp := value.(string)
		if _maxNum, err := strconv.Atoi(rule); err != nil {
			return fmt.Errorf("%s配置错误%s", name, rule)
		} else {
			if len([]rune(_tmp)) > _maxNum {
				return fmt.Errorf("%s长度不能大于%s", name, rule)
			}
		}
	}
	return nil
}

//最小 支持数字和字符长度
func validMin(value interface{}, rule, name string) error {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		_tmp := reflect.ValueOf(value).Int()
		if _Num, err := strconv.ParseInt(rule, 10, 64); err != nil {
			return fmt.Errorf("%s配置错误", rule)
		} else {
			if (_tmp < _Num) {
				return fmt.Errorf("%s不能小于%s", name, rule)
			}
		}
		break
	case uint, uint8, uint16, uint32, uint64:
		_tmp := reflect.ValueOf(value).Uint()
		if _Num, err := strconv.ParseUint(rule, 10, 64); err != nil {
			return fmt.Errorf("%s配置错误", rule)
		} else {
			if (_tmp > _Num) {
				return fmt.Errorf("%s不能小于%s", name, rule)
			}
		}

	case float32, float64:
		_tmp := reflect.ValueOf(value).Float()
		if _Num, err := strconv.ParseFloat(rule, 64); err != nil {
			return fmt.Errorf("%s配置错误", rule)
		} else {
			if (_tmp > _Num) {
				return fmt.Errorf("%s不能小于%s", name, rule)
			}
		}
	case string:
		_tmp := value.(string)
		if _Num, err := strconv.Atoi(rule); err != nil {
			return fmt.Errorf("%s配置错误%s", name, rule)
		} else {
			if len([]rune(_tmp)) < _Num {
				return fmt.Errorf("%s长度不能小于%s", name, rule)
			}
		}
	}
	return nil
}

/**
时间字符串验证
2006-01-02 15:04:05
 */
func validDate(date interface{}, rule string, name string) error {
	d, ok := date.(string)
	if !ok {
		return fmt.Errorf("%s类型不正确", name)
	}
	_, err := time.Parse(rule, d)
	if err != nil {
		return fmt.Errorf("%s格式不正确", name)
	}
	return nil
}

func validLen(value interface{}, rule string, name string) error {
	if d, ok := value.(string); ok {
		if _Num, err := strconv.Atoi(rule); err != nil {
			return fmt.Errorf("%s配置错误", name)
		} else {
			if len(d) != _Num {
				return fmt.Errorf("%s长度必须为%s", name, rule)
			}
		}
	}
	return nil
}

func validPhone(value interface{}, rule string, name string) error {
	d, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s类型不正确", name)
	}
	reg := regexp.MustCompile(`^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\d{8}$`)
	if reg.MatchString(d) == true {
		return nil
	} else {
		return fmt.Errorf("无效%s", name)
	}
}

func validEmail(value interface{}, rule, name string) error {
	d, ok:= value.(string)
	if !ok{
		return fmt.Errorf("%s类型不正确",name)
	}
	reg:=regexp.MustCompile(`^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`)
	if reg.MatchString(d) == true {
		return nil
	} else {
		return fmt.Errorf("无效%s", name)
	}
}
