package config

import (
	"testing"
	"fmt"
	"github.com/SunMaybo/jewel-template/template/hystrix"
	"encoding/json"
)

type Template struct {
	Service hystrix.Service `yaml:"service"`
}

func TestLoad(t *testing.T) {
	temp:=Template{}
	Load("test.yml", &temp)
	resp,_:=json.Marshal(temp)
	fmt.Println(string(resp))
}
