# go-http-client
go的http请求库参数特别多，有一些默认参数甚至会让你踩坑，我做了四组不同的实验来验证它

>server1 && server2
>> 服务端，监听一个端口，支持keepalive，并且不做超时限制（这是重点，真有这样的服务器）

>client1
>>全局变量发起请求，不做连接后空闲超时限制，请求连接上会一直存在，请求多个域名，每个域名有一个句柄。

>client2
>>全局变量发起，先不设置空闲时间请求，句柄为ESTAB, 然后设置空闲超时时间再请求，句柄变为TIME-WAIT 
>>全局变量的特点就是不同请求可能要求的超时等设置不同，他们会相互影响。只可在某一固定功能上通用

>client3
>>局部变量，不做连接后空闲超时限制，每个请求一个句柄，如果正好遇到服务器也没有空闲超时限制，那句柄数很快会爆

>client4
>>局部变量，设置空闲超时时间。每个请求一个句柄，句柄会正常关闭。不会有问题

# 参考
* https://blog.csdn.net/bigwhite20xx/article/details/112386441
* https://blog.huati365.com/ab420f833f7ffcfd
* https://wangchujiang.com/linux-command/c/ss.html