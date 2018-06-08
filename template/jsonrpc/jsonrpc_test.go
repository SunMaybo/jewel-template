package jsonrpc

import (
	"testing"
	"jewel-template/template/rest"
	"time"
	"fmt"
)

func TestClient_Call(t *testing.T) {
	client := Config(rest.ClientConfig{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
		SocketTimeout:      time.Second * 10,
		Authorization:      "Basic dG9rZW51cDp0b2tlbnVwLWJpdGNvaW4=",
	}, "http://13.251.28.17:8868")
	resp, _ := client.Call(RequestBody{
		Method: "listunspent",
		Id:     1,
		Params: []interface{}{0},
	})
	fmt.Println(fmt.Sprintf("%+v", resp))
}
