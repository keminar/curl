# go-http-client
go的http请求库参数特别多，有一些默认参数甚至会让你踩坑，我做了四组不同的实验来验证它

>server1 && server2
>> 服务端，监听一个端口，支持keepalive，并且不做超时限制（这是重点，真有这样的服务器）

>client1
>>全局变量发起请求，不做连接后空闲超时限制，请求连接上会一直存在，请求多个域名，每个域名有一个句柄。

>client2
>>全局变量发起，先不设置空闲时间请求，句柄为ESTAB, 然后设置IdleConnTimeout再请求，句柄变为TIME-WAIT 
>>全局变量的特点就是不同请求可能要求的超时等设置不同，他们会相互影响。只可在某一固定功能上通用

>client3
>>局部变量，不做连接后空闲超时限制，每个请求一个句柄，如果正好遇到服务器也没有空闲超时限制，那句柄数很快会爆

>client4
>>局部变量，设置SetDeadline。每个请求一个句柄，句柄会正常关闭。不会有问题

>client5
>>局部变量，关闭keepalive功能，请求的句柄使用完马上被销毁。观察返回值每次请求的端口都不同，说明都是重新TCP握手的

# 汇总

1. 测试中只要ts.IdleConnTimeout或设置conn.SetDeadline()都可以起到效果, 至于区别可看最下面的参考链接。
2. 根据client-timeout.png中所示，最好要设置3个超时时间（Dialer.Timeout,Client.Timeout,IdleConnTimeout)
3. 根据文中说明MaxIdleConnsPerHost默认值只有2有时也是不够的，可以设置为1024
4. 可以通过DisableKeepAlives:true关闭长连接

# 参考
* https://blog.csdn.net/bigwhite20xx/article/details/112386441
* https://blog.huati365.com/ab420f833f7ffcfd
* https://wangchujiang.com/linux-command/c/ss.html
* https://www.oschina.net/translate/the-complete-guide-to-golang-net-http-timeouts?print