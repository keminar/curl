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
	time.Sleep(time.Duration(20) * time.Second)
	if ts, ok := client.Transport.(*http.Transport); ok {
		// 这里是和client1唯一不同的地方
		ts.IdleConnTimeout = time.Duration(3) * time.Second //最大空闲时间
	}
	for i := 0; i < 10; i++ {
		send()
	}
	time.Sleep(time.Duration(1) * time.Hour)
}

var tr *http.Transport
var client *http.Client

func init() {
	tr = &http.Transport{
		DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, time.Duration(1)*time.Second)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(5) * time.Second, //单次请求超时时间
	}
}
func send() {
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
# ss查看句柄，只有一个状态为ESTAB的连接, 但当ts.IdleConnTimeout执行了以后ESTAB变和TIME-WAIT了

$ ss -t -a|grep 127.0.0.1:7788
ESTAB      0      0               127.0.0.1:37800               127.0.0.1:7788

$ ss -t -a|grep 127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37800           127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37802           127.0.0.1:7788
**/
