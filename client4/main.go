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
			// 这里是和client3唯一不同的地方
			conn.SetDeadline(time.Now().Add(time.Duration(3) * time.Second))
			return conn, nil
		},
	}
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
# ss查看句柄，有10个状态为TIME-WAIT的连接, 会正常断开链接，不容易出现句柄过多

$ ss -t -a|grep 127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37830           127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37852           127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37856           127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37828           127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37840           127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37858           127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37868           127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37860           127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37844           127.0.0.1:7788
TIME-WAIT  0      0          127.0.0.1:37864           127.0.0.1:7788

**/
