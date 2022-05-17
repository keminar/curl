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
# ss查看句柄，有10个状态为ESTAB的连接, 这种在客户端和服务端都没有空闲超时断开的设置下，文件打开句柄会越来越多，直到系统挂掉

$ ss -t -a|grep 127.0.0.1:7788
ESTAB      0      0               127.0.0.1:37550               127.0.0.1:7788
ESTAB      0      0               127.0.0.1:37552               127.0.0.1:7788
ESTAB      0      0               127.0.0.1:37546               127.0.0.1:7788
ESTAB      0      0               127.0.0.1:37554               127.0.0.1:7788
ESTAB      0      0               127.0.0.1:37536               127.0.0.1:7788
ESTAB      0      0               127.0.0.1:37534               127.0.0.1:7788
ESTAB      0      0               127.0.0.1:37548               127.0.0.1:7788
ESTAB      0      0               127.0.0.1:37542               127.0.0.1:7788
ESTAB      0      0               127.0.0.1:37540               127.0.0.1:7788
ESTAB      0      0               127.0.0.1:37544               127.0.0.1:7788

**/
