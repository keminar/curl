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
		time.Sleep(time.Duration(1) * time.Second)
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
			// 和client6不同之处，使用SetDeadline
			conn.SetDeadline(time.Now().Add(time.Duration(3) * time.Second))
			return conn, nil
		},
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(5) * time.Second, //单次请求超时时间
	}
}
func send() {
	apiUrl := "http://127.0.0.1:8877/"
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
			log.Println("read", string(body))
		}
	}
}

/**
# ss查看句柄，有四个状态为TIME-WAIT的连接

$ ss -t -a|grep 127.0.0.1:|grep 77
TIME-WAIT  0      0          127.0.0.1:37182           127.0.0.1:8877
TIME-WAIT  0      0          127.0.0.1:37184           127.0.0.1:8877
TIME-WAIT  0      0          127.0.0.1:37188           127.0.0.1:8877
TIME-WAIT  0      0          127.0.0.1:37186           127.0.0.1:8877

**/
