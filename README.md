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
		ReplyCount:3,
	})
	var resp map[string]interface{}
	err := restTemplate.GetForObject("http://13.229.115.157:8877/info", &resp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("%+v", resp))
```

### 熔断
```golang
    hystrixTable := make(hystrix.HystrixTable)
	hystrixTable["images"] = hystrix.Hystrix{Path: "/api/tasks/:id",RequestVolumeThreshold:9,RequestWindowsTime:30,ErrorPercentThreshold:35}
	serviceBucket := make(ServiceBucket)
	serviceBucket["storage_service"] = hystrix.Service{
		Host:           "192.168.1.100:31002",
		HystrixEnabled: true,
		HystrixTable:   hystrixTable,
	}
	factory := New(Config{
		JewelTemplate: JewelTemplate{
			Template: Template{
				ServiceBucket: serviceBucket,
			},
		},
	}, func(name string, isOpen bool) {
		fmt.Println(name,isOpen)
	})
    template := factory.Service("storage_service")
	dataMap := make(map[string]interface{})
	err := template.GetForObject("images", &dataMap, "a2ea2f3b771311e98f13a580af40044")
	if err != nil {
		fmt.Println(err.Status)
		log.Fatal(err)
	}
	fmt.Println(dataMap)
```
### 熔断配置更好结合
```yml
jewel:
   template:
       service:
         images_service:
             schema: https
             host: www.baidu.com
             hystrix_enabled: true
             rest:
               max_idle_conns: 5
               max_idle_timeout: 3000
               disable_compression: true
               socket_timeout: 3000
               reply_count: 3
               proxy: http://127.0.0.1:1087
             hystrix:
                 links:
                   path: /links
                   request_volume_threshold: 3
                   error_percent_threshold: 25
                   request_windows_time: 10
                 test:
                   path: /test
                   request_volume_threshold: 3
                   error_percent_threshold: 25
                   request_windows_time: 10
         article_service:
                      schema: http
                      host: 192.168.1.100:31002
                      hystrix_enabled: true
                      rest:
                        max_idle_conns: 5
                        max_idle_timeout: 3000
                        socket_timeout: 3000
                        reply_count: 3
                      hystrix:
                          links:
                            path: /api/tasks/:id
                            request_volume_threshold: 3
                            error_percent_threshold: 25
                            request_windows_time: 10
                          test:
                            path: /test
                            request_volume_threshold: 3
                            error_percent_threshold: 25
                            request_windows_time: 10
```
### 使用
```
filedata, err := readAll("test.yml")
	if err != nil {
		log.Fatal(err)
	}
	config := template.Config{}
	err = yaml.Unmarshal(filedata, &config)
	if err != nil {
		log.Fatal(err)
	}
	buff, _ := json.Marshal(config)
	fmt.Println(string(buff))
	factory = template.New(config, func(name string, isOpen bool) {
		fmt.Println(name, isOpen)
	})
```
```配置参数描述
jewel:
   template:
       service:
         images_service:                                     服务名字
             schema: https                                   请求模式http or https
             host: www.baidu.com                             域名
             hystrix_enabled: true                           是否熔断
             rest:                          httpclient 配置
               max_idle_conns: 5                             最大闲置连接数
               idle_conn_timeout: 3000                       闲置连接超时时间ms
               disable_compression: true                     是否压缩
               socket_timeout: 3000                          请求时间ms
               reply_count: 3                                重试次数
               proxy: http://127.0.0.1:1087                  代理
             hystrix:                          熔断配置
                 links:                                      请求名字
                   path: /links                              路径
                   request_volume_threshold: 3               时间窗口最小请求数
                   error_percent_threshold: 25               失败率0~100
                   request_windows_time: 10                  时间窗口大小s
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

    