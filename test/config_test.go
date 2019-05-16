package test

import (
	"os"
	"log"
	"jewel-template/template"
	"gopkg.in/yaml.v2"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"testing"
)

var factory = template.JewelTemplateFactory{}

func init() {
	filedata, err := readAll("test.yml")
	if err != nil {
		log.Fatal(err)
	}
	config := template.Config{}
	err = yaml.Unmarshal(filedata, &config)
	if err != nil {
		log.Fatal(err)
	}
	buff, _ := json.Marshal(config)
	fmt.Println(string(buff))
	factory = template.New(config, func(name string, isOpen bool) {
		fmt.Println(name, isOpen)
	})
}

func readAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}
func TestGetHystrixTemplate(t *testing.T) {
	template := factory.Service("article_service")
	dataMap := make(map[string]interface{})
	err := template.Execute("links", "GET", nil, nil, &dataMap, "a2ea2f3b771311e98f130a580af40044")
	if err != nil {
		log.Fatal(err)
	}
}
