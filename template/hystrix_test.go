package template

import (
	"testing"
	"github.com/SunMaybo/jewel-template/template/hystrix"
	"github.com/bdlm/log"
	"fmt"
	"errors"
)

var factory JewelTemplateFactory

func init() {
	hystrixTable := make(hystrix.HystrixTable)
	hystrixTable["images"] = hystrix.Hystrix{Path: "/api/tasks/:id"}
	serviceBucket := make(ServiceBucket)
	serviceBucket["storage_service"] = hystrix.Service{
		Host:           "192.168.1.100:31002",
		HystrixEnabled: true,
		HystrixTable:   hystrixTable,
	}
	factory = New(Config{
		JewelTemplate: JewelTemplate{
			Template: Template{
				ServiceBucket: serviceBucket,
			},
		},
	})
}

func TestHystrix(t *testing.T) {
	template := factory.Service("storage_service", nil)
	dataMap := make(map[string]interface{})
	err := template.GetForObject("images", &dataMap, "a2ea2f3b771311e98f13a580af40044")
	if err != nil {
		fmt.Println(err.Status)
		log.Fatal(err)
	}
	fmt.Println(dataMap)
}
func TestResponseHystrix(t *testing.T) {
	template := factory.Service("storage_service")
	dataMap := make(map[string]interface{})
	err := template.ExecuteWithCustomHystrix("images", "GET", nil, nil, &dataMap, func(response interface{}) error {
		data := response.(*map[string]interface{})
		if (*data)["status"] == "FINISHED" {
			return errors.New("errr")
		}

		return nil
	}, "a2ea2f3b771311e98f130a580af40044")
	if err != nil {
		fmt.Println(err.Status)
		log.Fatal(err)
	}
}
