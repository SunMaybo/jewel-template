package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"reflect"
	"strings"
	"os"
	"strconv"
	"github.com/bdlm/log"
)

func Load(fileName string, out interface{}) {
	buff, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(buff, out)
	if err != nil {
		log.Fatal(err)
	}
	loadOSEnv(out)
}
func loadOSEnv(out interface{}) {
	outType := reflect.TypeOf(out)
	outValue := reflect.ValueOf(out)
	scan(outType, outValue, "")
}

func scan(outType reflect.Type, outValue reflect.Value, prefix string) {
	//指针处理
	if outType.Kind() == reflect.Ptr { //指针类型处理
		scan(outType.Elem(), outValue.Elem(), prefix)
	}
	//结构体类型
	if outType.Kind() == reflect.Struct {
		for i := 0; i < outType.NumField(); i++ {
			sField := outType.Field(i)
			if sField.Tag.Get("yaml") == "" {
				continue
			}
			sValue := outValue.Field(i)
			if prefix == "" {
				prefix = sField.Tag.Get("yaml")
			} else {
				prefix = prefix + "." + sField.Tag.Get("yaml")
			}
			scan(sField.Type, sValue, prefix)
		}
	}

	if prefix == "" {
		return
	}
	osEnvStr := os.Getenv(strings.ToUpper(prefix))
	if osEnvStr == "" {
		return
	}
	//整形处理
	if outType.Kind() >= reflect.Int && outType.Kind() <= reflect.Uint64 {
		osEnv, err := strconv.ParseInt(osEnvStr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		outValue.SetInt(osEnv)
	}
	//布尔值处理
	if outType.Kind() == reflect.Bool {
		osEnv, err := strconv.ParseBool(osEnvStr)
		if err != nil {
			log.Fatal(err)
		}
		outValue.SetBool(osEnv)
	}
	//浮点值处理
	if outType.Kind() >= reflect.Float32 && outType.Kind() <= reflect.Float64 {
		osEnv, err := strconv.ParseFloat(osEnvStr, 64)
		if err != nil {
			log.Fatal(err)
		}
		outValue.SetFloat(osEnv)
	}
	//字符串处理
	if outType.Kind() >= reflect.String {
		outValue.SetString(osEnvStr)
	}
	return
}
