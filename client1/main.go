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
		send(i)
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
func send(i int) {
	apiUrl := "http://127.0.0.1:8877/"
	if i%2 == 0 {
		apiUrl = "http://127.0.0.1:7788/"
	}
	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client do:", err)
	} else {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			log.Println("read", apiUrl, string(body))
		}
	}
}

/**
# ss查看句柄，每个域名只有一个状态为ESTAB的连接

$ ss -t -a|grep 127.0.0.1:|grep 77
ESTAB      0      0                 127.0.0.1:37144               127.0.0.1:8877
ESTAB      0      0                 127.0.0.1:37948               127.0.0.1:7788

**/
