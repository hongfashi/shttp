package stools

import(
	"os"
	"net"
	"net/http"
	"net/url"
	"compress/gzip"
	"bytes"
	"io"
	"io/ioutil"
	"fmt"
	"time"
	"encoding/json"
	"strings"
	
)

func HttpDo(client http.Client,method string,url string,headers map[string] string,reader io.Reader)([]byte){
 	
	content := []byte(`{}`)
	
	// 使用Get方法获取服务器响应包数据
   req, err := http.NewRequest(method,url,reader)
   
   if err != nil {
      fmt.Println("request err:", err)
      return content
   }
    for k, v := range headers {
		req.Header.Add(k,v);
        //fmt.Println("k,v = ", k,v) 
    }
	
	
	resp, err := client.Do(req)
   
   if err != nil {
      fmt.Println("client err:", err)
      return content
   }
   
   defer resp.Body.Close()
	
	// 获取服务器端读到的数据
	
   //fmt.Println("Status = ", resp.Status)           // 状态
   //fmt.Println("StatusCode = ", resp.StatusCode)   // 状态码
   //fmt.Println("Header = ", resp.Header)           //响应头
   //fmt.Println("Body = ", resp.Body)               // 响应包体
	//读取body内的内容
	
	
	body, err := gzip.NewReader(resp.Body)
   
    if err != nil {
      fmt.Println("gzip err:", err)
      return content
   }
	defer body.Close()
	
	content, err = ioutil.ReadAll(body)
	
	//fmt.Println("result",string(content))
	
	 if err != nil {
      fmt.Println("io read err:", err)
      return content
   }
   
   return content
 
}

func HttpJson(client http.Client,method string,url string,headers map[string] string,postData map[string] string)([]byte){
	
	content := []byte(`{}`)
	
	data,err := json.Marshal(postData)
	
	//strlen := bytes.Count(data,nil)-1
	//fmt.Println(strlen)
	//return  content
	
	if err!=nil {
	
		fmt.Println(err)
		return content
	}
	
	reader := bytes.NewReader(data)
	
	return HttpDo(client,method,url,headers,reader)
	
}

func HttpPost(client http.Client,method string,url string,headers map[string] string,postData map[string] string)([]byte){
	
	
	post := FormData(postData)
	
	
	//fmt.Println("post=",post)
	reader := strings.NewReader(post)
	
	return HttpDo(client,method,url,headers,reader)
	
}



func NewHttpClient(proxyAddr string) http.Client {
    proxy, _ := url.Parse(proxyAddr)
  
    netTransport := &http.Transport{
        //Proxy: http.ProxyFromEnvironment,
        Proxy: http.ProxyURL(proxy),
        Dial: func(netw, addr string) (net.Conn, error) {
            c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(10))
            if err != nil {
                return nil, err
            }
            return c, nil
        },
        MaxIdleConnsPerHost:   10,                            //每个host最大空闲连接
        ResponseHeaderTimeout: time.Second * time.Duration(5), //数据收发5秒超时
    }

    return http.Client{
        Timeout:   time.Second * 10,
        Transport: netTransport,
    }
}
