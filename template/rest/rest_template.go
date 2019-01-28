package rest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"fmt"
	"errors"
	"regexp"
	"encoding/json"
	"time"
	"log"
	"net/url"
)

type Template struct {
	Client      *http.Client
	EnableReply bool
	ReplyCount  int
}

type ClientConfig struct {
	MaxIdleConns       int
	IdleConnTimeout    time.Duration
	DisableCompression bool
	SocketTimeout      time.Duration
	Authorization      string
	ReplyCount         int
	Proxy              string
}

func Default() *RestTemplate {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    3 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	client.Timeout = 10 * time.Second
	return &RestTemplate{
		Template: Template{
			Client:      client,
			EnableReply: true,
		},
	}
}
func DefaultProxy(proxy string) *RestTemplate {
	u, _ := url.Parse(proxy)
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    3 * time.Second,
		DisableCompression: true,
		Proxy:              http.ProxyURL(u),
	}
	client := &http.Client{Transport: tr}
	client.Timeout = 10 * time.Second
	return &RestTemplate{
		Template: Template{
			Client:      client,
			EnableReply: true,
		},
	}
}

func Config(cfg ClientConfig) *RestTemplate {
	var tr *http.Transport
	if cfg.Proxy != "" {
		u, _ := url.Parse(cfg.Proxy)
		tr = &http.Transport{
			MaxIdleConns:       cfg.MaxIdleConns,
			IdleConnTimeout:    cfg.IdleConnTimeout,
			DisableCompression: cfg.DisableCompression,
			Proxy:              http.ProxyURL(u),
		}

	} else {
		tr = &http.Transport{
			MaxIdleConns:       cfg.MaxIdleConns,
			IdleConnTimeout:    cfg.IdleConnTimeout,
			DisableCompression: cfg.DisableCompression,
		}
	}

	client := &http.Client{Transport: tr}
	client.Timeout = cfg.SocketTimeout
	return &RestTemplate{
		Template: Template{
			Client:      client,
			ReplyCount:  cfg.ReplyCount,
			EnableReply: true,
		},
		ClientConfig: &cfg,
	}
}

func (template *Template) Call(url string, param []byte, method string, Header http.Header) ([]byte, error) {
	reader := bytes.NewReader(param)
	var req *http.Request
	var err error
	req, err = http.NewRequest(method, url, reader)
	if err != nil {
		return []byte{}, err
	}
	req.Close = true
	if Header != nil {
		req.Header = Header
	}
	resp, err := template.Client.Do(req)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	BodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return BodyByte, errors.New(fmt.Sprintf("%d-%s-%s", resp.StatusCode, "error", string(BodyByte)))
	}
	return BodyByte, nil
}

func (template *Template) CallWithReply(url string, param []byte, method string, header http.Header, count int) ([]byte, error) {
	reader := bytes.NewReader(param)
	var req *http.Request
	var err error
	req, err = http.NewRequest(method, url, reader)
	if err != nil {
		return []byte{}, err
	}
	req.Close = true
	if header != nil {
		req.Header = header
	}
	resp, err := template.Client.Do(req)
	if err != nil && count < template.ReplyCount {
		count++
		return template.CallWithReply(url, param, method, header, count)
	} else if err != nil && count >= template.ReplyCount {
		return []byte{}, err
	}

	BodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return BodyByte, errors.New(fmt.Sprintf("%d-%s-%s", resp.StatusCode, "error", string(BodyByte)))
	}
	return BodyByte, nil
}

type RestTemplate struct {
	Template
	ClientConfig *ClientConfig
}

func (rest *RestTemplate) GetForObject(url string, response interface{}, uriVariables ... string) error {
	return rest.ExecuteForObject(url, http.MethodGet, nil, nil, response, uriVariables...)
}

func (rest *RestTemplate) PostForObject(url string, body, response interface{}, uriVariables ... string) error {
	return rest.ExecuteForObject(url, http.MethodPost, nil, body, response, uriVariables...)
}
func (rest *RestTemplate) PutForObject(url string, body, response interface{}, uriVariables ... string) error {
	return rest.ExecuteForObject(url, http.MethodPut, nil, body, response, uriVariables...)
}

func (rest *RestTemplate) DeleteForObject(url string, response interface{}, uriVariables ... string) error {
	return rest.ExecuteForObject(url, http.MethodDelete, nil, nil, response, uriVariables...)
}

func (rest *RestTemplate) HeadForObject(url string, header http.Header, response interface{}, uriVariables ... string) error {
	return rest.ExecuteForObject(url, http.MethodHead, header, nil, response, uriVariables...)
}
func (rest *RestTemplate) ExecuteForJsonString(url, method string, header http.Header, body string, response interface{}, uriVariables ... string) error {
	url = convertToUrl(url, uriVariables...)
	if header == nil {
		header = http.Header{}
	}
	header.Set("Content-Type", "application/json")
	if rest.ClientConfig != nil && rest.ClientConfig.Authorization != "" {
		header.Set("Authorization", rest.ClientConfig.Authorization)
	}
	buff := []byte(body)
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if rest.EnableReply {
		result, err := rest.CallWithReply(url, buff, method, header, 0)
		if err != nil {
			return err
		}
		json.Unmarshal(result, response)
	} else {
		result, err := rest.Call(url, buff, method, header)
		if err != nil {
			return err
		}
		json.Unmarshal(result, response)
	}
	return nil
}

func (rest *RestTemplate) ExecuteForObject(url, method string, header http.Header, body, response interface{}, uriVariables ... string) error {
	url = convertToUrl(url, uriVariables...)
	if header == nil {
		header = http.Header{}
	}
	header.Set("Content-Type", "application/json")
	if rest.ClientConfig != nil && rest.ClientConfig.Authorization != "" {
		header.Set("Authorization", rest.ClientConfig.Authorization)
	}
	buff, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if rest.EnableReply {
		result, err := rest.CallWithReply(url, buff, method, header, 0)
		if err != nil {
			return err
		}
		json.Unmarshal(result, response)
	} else {
		result, err := rest.Call(url, buff, method, header)
		if err != nil {
			return err
		}
		json.Unmarshal(result, response)
	}
	return nil
}

func (rest *RestTemplate) Execute(url, method string, header http.Header, body, response interface{}, uriVariables ... string) error {
	url = convertToUrl(url, uriVariables...)
	if header == nil {
		header = http.Header{}
	}
	if rest.ClientConfig != nil && rest.ClientConfig.Authorization != "" {
		header.Set("Authorization", rest.ClientConfig.Authorization)
	}

	buff, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if rest.EnableReply {
		result, err := rest.CallWithReply(url, buff, method, header, 0)
		if err != nil {
			return err
		}
		json.Unmarshal(result, response)
	} else {
		result, err := rest.Call(url, buff, method, header)
		if err != nil {
			return err
		}
		json.Unmarshal(result, response)
	}
	return nil
}

func convertToUrl(url string, uriVariables ... string) string {
	regex := `:([a-zA-Z_][^/?&]+)[/?&]{0,1}`
	re, _ := regexp.Compile(regex)
	var path []byte
	urlByte := []byte(url)
	all := re.FindAllSubmatchIndex(urlByte, -1)
	for i, point := range all {
		if i == 0 {
			path = append(path, urlByte[:point[2]-1]...)
		} else {
			path = append(path, urlByte[(all[i-1])[3]:point[2]-1]...)
		}
		if len(uriVariables) <= i {
			path = append(path, urlByte[(all[i])[2]-1:point[3]]...)
		} else {
			path = append(path, []byte(uriVariables[i])...)
		}

	}
	if len(all) == 0 {
		path = append(path, urlByte[:]...)
	} else {
		path = append(path, urlByte[all[len(all)-1][3]:]...)
	}
	return string(path)
}
