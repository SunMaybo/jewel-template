package jsonrpc

import (
	"github.com/SunMaybo/jewel-template/rest"
	"net/http"
	"time"
)

type RequestBody struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int64         `json:"id"`
}
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type ResponseBody struct {
	JsonRpc string      `json:"jsonrpc"`
	Error   error       `json:"error,omitempty"`
	Result  interface{} `json:"result"`
	Id      int64       `json:"id"`
}
type Client struct {
	Url      string
	template rest.RestTemplate
}

func Config(cfg rest.ClientConfig, url string) *Client {
	tr := &http.Transport{
		MaxIdleConns:       cfg.MaxIdleConns,
		IdleConnTimeout:    cfg.IdleConnTimeout,
		DisableCompression: cfg.DisableCompression,
	}
	client := &http.Client{Transport: tr}
	client.Timeout = cfg.SocketTimeout
	return &Client{
		template: rest.RestTemplate{
			Template: rest.Template{
				Client: client,
			},
			ClientConfig: &cfg,
		},
		Url: url,
	}
}
func Default() *Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    3 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	client.Timeout = 5 * time.Second
	return &Client{
		template: rest.RestTemplate{
			Template: rest.Template{
				Client: client,
			},
		},
		Url: "http://127.0.0.1:8080",
	}
}
func (client *Client) Call(body RequestBody) (ResponseBody, error) {
	response := ResponseBody{}
	body.JsonRpc = "2.0"
	return response, client.template.PostForObject(client.Url, body, &response)
}
