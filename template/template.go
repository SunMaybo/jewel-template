package template

import "jewel-template/template/hystrix"

type ServiceBucket map[string]hystrix.Service
type Template struct {
	ServiceBucket ServiceBucket `yaml:"Service"`
}
type JewelTemplate struct {
	Template Template `yaml:"template"`
}

type Config struct {
	JewelTemplate JewelTemplate `yaml:"jewel"`
}
type JewelTemplateFactory struct {
	config Config
}

func New(config Config) JewelTemplateFactory {
	return JewelTemplateFactory{config: config}
}
func (jtf JewelTemplateFactory) Service(name string) *hystrix.HystrixTemplate {
	if service, ok := jtf.config.JewelTemplate.Template.ServiceBucket[name]; ok {
		return hystrix.New(service)
	}
	return nil
}
