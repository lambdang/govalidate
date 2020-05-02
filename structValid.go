package govalidate

import (
	"fmt"
	"reflect"
	"strings"
)



func ValidStruct(stc interface{},validEn ...ValidStore)(map[string]string,bool){
	fields := reflect.TypeOf(stc).Elem()
	values := reflect.ValueOf(stc).Elem()
	_errors := map[string]string{}
	isok := true
	var validManager ValidStore
	if len(validEn)!=0{
		validManager = validEn[0]
	}else{
		validManager = NewValidStore()

	}
	for i := 0; i < fields.NumField(); i++ {

		field := fields.Field(i)
		tag := field.Tag
		strName := tag.Get("name") //有name 则用name名字接收
		validRule :=tag.Get("valid")
		jsonTag := tag.Get("json")

		fieldName := field.Name //struct的name
		JsonName := strings.Split(jsonTag,",")[0]
		if validRule!=""{ //需要验证
			mapRule,err:=ParseRule(validRule)
			if err!=nil{
				es:=fmt.Sprintf("%s  %s",JsonName,err.Error())
				_errors["err"] = es
				return _errors,false
			}else{
				value := values.Field(i)
				if strName == ""{
					strName = fieldName
				}
				err,ok:=validManager.ValidRule(value,strName,mapRule) //只抛出单条异常
				if !ok {
					isok=false
					_errors[JsonName] = err.Error()
				}
			}

		}
	}

	return _errors,isok
}


