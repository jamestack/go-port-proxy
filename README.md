# go-port-proxy
这是一个用Go语言写的端口转发工具，支持TCP、UDP。
## 编译与运行
```
#编译
git clone github.com/jameswone/go-port-proxy.git
cd go-port-proxy/
go build port-proxy.go
#运行
./port-proxy tcp 127.0.0.1:80 192.168.1.102:8080
```
# tcp-websocket-proxy
这是一个用Go语言写的安全通道流量转发工具，支持通过Websocket通道转发端口流量。
## 用途
可以用于网络防火墙穿透！CROSS GFW!
例如：利用CDN翻墙，防止Shadowsocks服务器ip被封！<br/>
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;客户端(Shadowsocks+Proxy)&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;服务器(Shadowsocks+Proxy)<br/>
SS-client <-tcp-> Proxy-client <-ws-> Cdn <-ws-> Proxy-server <-tcp-> SS-server
## 编译与运行
```
#编译
git clone github.com/jameswone/go-port-proxy.git
cd go-port-proxy/
go build tcp-ws-proxy.go
#运行
./tcp-ws-proxy tcp-to-ws 127.0.0.1:8087 www.jamestack.tk:8080
```
## 利用Cdn代理Shadowsocks流量
```
#客户端：
./tcp-ws-proxy tcp-to-ws 0.0.0.0:8087 www.jamestack.tk:8080
#服务器：
./tcp-ws-proxy ws-to-tcp 0.0.0.0:8080 0.0.0.0:8087
```
