package rest

import (
	"testing"
	"fmt"
	"encoding/json"
	"log"
)

func TestUriRegexp(t *testing.T) {
	data := "http://192.168.200.119:8080/test009_90/name/:age"
	fmt.Println(json.Marshal(nil))
	fmt.Println(convertToUrl(data, "xxxxxx"))
}
func TestGet(t *testing.T) {
	restTemplate := Default()
	var resp map[string]interface{}
	err := restTemplate.GetForObject("http://13.229.115.157:8877/info", &resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("%+v", resp))
}

