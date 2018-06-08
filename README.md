# jewel-template
基于JAVA RestTemplate 思想实现对golang http.client 封装，方便模版使用
## 适用于
    1. 基于jsonrpc2.0 接口调用
    2. 基于restful风格的接口
### 适用实例(rest)
#### 默认配置
```golang
restTemplate := rest.Default()
err := restTemplate.GetForObject("http://127.0.0.1:8877/info", &resp)
	if err != nil {
		log.Fatal(err)
	}
fmt.Println(fmt.Sprintf("%+v", resp))
```
#### 自定义配置
```golang
	restTemplate := rest.Config(rest.ClientConfig{
		MaxIdleConns:20,
		IdleConnTimeout:3*time.Second,
		DisableCompression:true,
		SocketTimeout:3*time.Second,
	})
	var resp map[string]interface{}
	err := restTemplate.GetForObject("http://13.229.115.157:8877/info", &resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("%+v", resp))
```
### 适用实例(jsonrpc2.0)
```golang
client := jsonrpc.Config(rest.ClientConfig{
		MaxIdleConns:       10,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
		SocketTimeout:      time.Second * 10,
		Authorization:      "Basic dG9rZW51cDp0b2tlbnVwLWJpdGNvaW4=",
	}, "http://127.0.0.1:8868")
	resp, _ := client.Call(jsonrpc.RequestBody{
		Method: "listunspent",
		Id:     1,
		Params: []interface{}{0},
	})
	fmt.Println(fmt.Sprintf("%+v", resp))

```

    