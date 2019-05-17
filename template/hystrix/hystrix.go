package hystrix

import (
	"net/http"
	"github.com/SunMaybo/jewel-template/template/rest"
	"github.com/SunMaybo/jewel-template/template/errors"
	"github.com/SunMaybo/hystrix-go/hystrix"
	"time"
)

//熔断配置
type Hystrix struct {
	Name                   string `yaml:"name"`
	Path                   string `yaml:"path"`
	RequestVolumeThreshold int    `yaml:"request_volume_threshold"`
	ErrorPercentThreshold  int    `yaml:"error_percent_threshold"`
	RequestWindowsTime     int64  `yaml:"request_windows_time"`
	//RequestVolumeThreshold uint64        `yaml:"request_volume_threshold"`
	//SleepWindow            time.Duration `yaml:"sleep_window"`
	//Timeout                time.Duration `yaml:"timeout"`
	//IsAlerting             bool   `yaml:"is_alerting"`
}
type HystrixTable map[string]Hystrix
type Service struct {
	Schema          string       `yaml:"schema"` //http or https and default http
	Name            string       `yaml:"name"`
	Host            string       `yaml:"host"`
	ClientConfig    *RestConfig  `yaml:"rest"`
	HystrixTable    HystrixTable `yaml:"hystrix"`
	HystrixEnabled  bool         `yaml:"hystrix_enabled"`
	hystrixTemplate *HystrixTemplate
}
type RestConfig struct {
	MaxIdleConns       int    `yaml:"max_idle_conns"`
	IdleConnTimeout    int    `yaml:"idle_conn_timeout"`
	DisableCompression bool   `yaml:"disable_compression"`
	SocketTimeout      int    `yaml:"socket_timeout"`
	Authorization      string `yaml:"authorization"`
	ReplyCount         int    `yaml:"reply_count"`
	Proxy              string `yaml:"proxy"`
}

type HystrixTemplate struct {
	rest    *rest.RestTemplate
	service Service
}

func New(service Service, hystrixFunc func(name string, isOpen bool)) (*HystrixTemplate) {
	if service.Schema == "" {
		service.Schema = "http"
	}
	ht := HystrixTemplate{
		service: service,
	}
	if service.ClientConfig != nil {
		ht.rest = rest.Config(rest.ClientConfig{
			MaxIdleConns:       service.ClientConfig.MaxIdleConns,
			IdleConnTimeout:    time.Millisecond * time.Duration(service.ClientConfig.IdleConnTimeout),
			DisableCompression: service.ClientConfig.DisableCompression,
			SocketTimeout:      time.Millisecond * time.Duration(service.ClientConfig.SocketTimeout),
			Authorization:      service.ClientConfig.Authorization,
			ReplyCount:         service.ClientConfig.ReplyCount,
			Proxy:              service.ClientConfig.Proxy,
		})
	} else {
		ht.rest = rest.Default()
	}

	for name, htx := range service.HystrixTable {
		htx.Name = name
		hystrix.ConfigureCommand(name, hystrix.CommandConfig{
			RequestVolumeThreshold: htx.RequestVolumeThreshold,
			ErrorPercentThreshold:  htx.ErrorPercentThreshold,
			RequestWindowsTime:     htx.RequestWindowsTime,
			Timeout:                60000,
			AlertFunc:              hystrixFunc,
		})
	}
	return &ht
}

func (t *HystrixTemplate) GetForObject(name string, response interface{}, uriVariables ... string) *errors.HttpError {
	return t.ExecuteForObject(name, http.MethodGet, nil, nil, response, uriVariables...)
}
func (t *HystrixTemplate) PostForObject(name string, body, response interface{}, uriVariables ... string) *errors.HttpError {
	return t.ExecuteForObject(name, http.MethodPost, nil, body, response, uriVariables...)
}
func (t *HystrixTemplate) PutForObject(name string, body, response interface{}, uriVariables ... string) *errors.HttpError {

	return t.ExecuteForObject(name, http.MethodPut, nil, body, response, uriVariables...)
}
func (t *HystrixTemplate) DeleteForObject(name string, response interface{}, uriVariables ... string) *errors.HttpError {

	return t.ExecuteForObject(name, http.MethodDelete, nil, nil, response, uriVariables...)
}
func (t *HystrixTemplate) HeadForObject(name string, header http.Header, response interface{}, uriVariables ... string) *errors.HttpError {

	return t.ExecuteForObject(name, http.MethodHead, header, nil, response, uriVariables...)
}
func (t *HystrixTemplate) ExecuteForJsonString(name, method string, header http.Header, body string, response interface{}, uriVariables ... string) *errors.HttpError {
	url, err := t.getUrl(name)
	if err != nil {
		return err
	}
	if !t.service.HystrixEnabled {
		return t.rest.ExecuteForJsonString(url, method, header, body, response, uriVariables...).(*errors.HttpError)
	}
	okChan := make(chan bool)
	errorChan := hystrix.Go(name, func() error {
		err := t.rest.ExecuteForJsonString(url, method, header, body, response, uriVariables...)
		if err == nil {
			okChan <- true
		}
		return err
	}, func(e error) error {
		return e
	})
	select {
	case err := <-errorChan:
		return errors.New(3005, err.Error()).(*errors.HttpError)
	case <-okChan:
		return nil
	}
}
func (t *HystrixTemplate) ExecuteForObject(name, method string, header http.Header, body, response interface{}, uriVariables ... string) *errors.HttpError {
	url, err := t.getUrl(name)
	if err != nil {
		return err
	}
	if !t.service.HystrixEnabled {
		return t.rest.ExecuteForObject(url, method, header, body, response, uriVariables...).(*errors.HttpError)
	}
	okChan := make(chan bool)
	errorChan := hystrix.Go(name, func() error {
		err := t.rest.ExecuteForObject(url, method, header, body, response, uriVariables...)
		if err == nil {
			okChan <- true
		}
		return err
	}, func(e error) error {
		return e
	})
	select {
	case err := <-errorChan:
		return errors.New(3005, err.Error()).(*errors.HttpError)
	case <-okChan:
		return nil
	}
}
func (t *HystrixTemplate) Execute(name, method string, header http.Header, body, response interface{}, uriVariables ... string) *errors.HttpError {
	url, err := t.getUrl(name)
	if err != nil {
		return err
	}
	if !t.service.HystrixEnabled {
		return t.rest.Execute(url, method, header, body, response, uriVariables...).(*errors.HttpError)
	}
	okChan := make(chan bool)
	errorChan := hystrix.Go(name, func() error {
		err := t.rest.Execute(url, method, header, body, response, uriVariables...)
		if err == nil {
			okChan <- true
		}
		return err
	}, func(e error) error {
		return e
	})
	select {
	case err := <-errorChan:
		return errors.New(3005, err.Error()).(*errors.HttpError)
	case <-okChan:
		return nil
	}
}
func (t *HystrixTemplate) ExecuteWithCustomHystrix(name, method string, header http.Header, body, response interface{}, responseCutBrokenFunc func(response interface{}) error, uriVariables ... string) *errors.HttpError {
	url, err := t.getUrl(name)
	if err != nil {
		return err
	}
	if !t.service.HystrixEnabled {
		return t.rest.Execute(url, method, header, body, response, uriVariables...).(*errors.HttpError)
	}
	okChan := make(chan bool)
	errorChan := hystrix.Go(name, func() error {
		err := t.rest.Execute(url, method, header, body, response, uriVariables...)
		if responseCutBrokenFunc != nil && err == nil {
			err = responseCutBrokenFunc(response)
		}
		if err == nil {
			okChan <- true
		}
		return err

	}, func(e error) error {
		return e
	})
	select {
	case err := <-errorChan:
		return errors.New(3005, err.Error()).(*errors.HttpError)

	case <-okChan:
		return nil
	}
}
func (t *HystrixTemplate) getUrl(name string) (string, *errors.HttpError) {
	if hystrix, ok := t.service.HystrixTable[name]; ok {
		return t.service.Schema + "://" + t.service.Host + hystrix.Path, nil
	}
	return "", errors.New(30006, "path is required").(*errors.HttpError)

}
