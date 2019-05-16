package template

import "github.com/SunMaybo/jewel-template/template/hystrix"

type ServiceBucket map[string]hystrix.Service
type Template struct {
	ServiceBucket ServiceBucket `yaml:"service"`
}
type JewelTemplate struct {
	Template Template `yaml:"template"`
}

type Config struct {
	JewelTemplate JewelTemplate `yaml:"jewel"`
}
type HystrixFunc func(name string, isOpen bool)

type JewelTemplateFactory struct {
	config      Config
	hystrixFunc HystrixFunc
}

func New(config Config, hystrixFunc HystrixFunc) JewelTemplateFactory {
	return JewelTemplateFactory{config: config, hystrixFunc: hystrixFunc}
}
func (jtf JewelTemplateFactory) Service(name string) *hystrix.HystrixTemplate {
	if service, ok := jtf.config.JewelTemplate.Template.ServiceBucket[name]; ok {
		return hystrix.New(service, jtf.hystrixFunc)
	}
	return nil
}
