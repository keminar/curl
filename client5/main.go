package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		send()
	}
	time.Sleep(time.Duration(1) * time.Hour)
}

func send() {
	tr := &http.Transport{
		DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, time.Duration(1)*time.Second)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
		//这是和client3不同之处
		DisableKeepAlives: true,
	}
	// 和client1的主要区别是变成了局部变量
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(5) * time.Second, //单次请求超时时间
	}
	apiUrl := "http://127.0.0.1:7788/"
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client do:", err)
	} else {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Println("read", string(body))
		}
	}
}

/**
# ss查看句柄，有0个句柄，因为瞬间就会被关闭了
**/
