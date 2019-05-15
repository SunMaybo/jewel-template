package template

import (
	"testing"
	"github.com/SunMaybo/jewel-template/template/hystrix"
	"github.com/bdlm/log"
	"fmt"
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

	template := factory.Service("storage_service")
	dataMap := make(map[string]interface{})
	err := template.GetForObject("images", dataMap, "65986c1d770c11e9a7cb0a580af4007a")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dataMap)
}
