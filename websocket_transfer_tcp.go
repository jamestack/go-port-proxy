package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"regexp"
)

func main() {

	argv := os.Args
	if len(argv) < 2 {
		fmt.Println("mode  localport WebSocketAddr")
		fmt.Println("ws-to-tcp 0.0.0.0:8088 0.0.0.0:8087")
		fmt.Println("tcp-to-ws 0.0.0.0:8080 www.jamestack.tk/asdasd:8087")
		return
	}

	listen, err := net.Listen("tcp", argv[2])
	if err != nil {
		fmt.Println("listen err.")
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("conn err.")
			break
		}
		fmt.Println("New Connect:", conn.RemoteAddr())
		switch argv[1] {
		case "ws-to-tcp":
			go onWsAccept(conn, argv[3])
		case "tcp-to-ws":
			go onTcpAccept(conn, argv[3])
		default:
			fmt.Println("argv err")
		}
	}

}

//WebSocket To Tcp
func onWsAccept(conn net.Conn, addr string) {
	//Websocket hand
	buff := make([]byte, 2048)
	num, err := conn.Read(buff)
	if err != nil {
		fmt.Println("read err.")
	}
	response := string(buff)
	fmt.Println(num, response)
	headder := make(map[string]string)
	if headpath := regexp.MustCompile(`(?:^)([^\r\n]+)(?:\r\n)`).FindAllStringSubmatch(response, -1); len(headpath) > 0 {
		headder["HeadPath"] = headpath[0][1]
	}
	if path := regexp.MustCompile(`(?:GET\s)(.+)(?:\sHTTP\/\d)`).FindAllStringSubmatch(response, -1); len(path) > 0 {
		headder["path"] = path[0][1]
	}
	headtitle := []string{"Host", "Connection", "Pragma", "Cache-Control", "Upgrade", "Origin", "Sec-WebSocket-Version", "DNT", "User-Agent", "Accept-Language", "Sec-WebSocket-Key", "Sec-WebSocket-Extensions"}
	for _, title := range headtitle {
		if reg := regexp.MustCompile(`(?:`+title+`:\s)(.+)(?:\r\n)`).FindAllStringSubmatch(response, -1); len(reg) > 0 {
			headder[title] = reg[0][1]
		}
	}
	if _, ok := headder["Sec-WebSocket-Key"]; ok {
		sha := sha1.New()
		sha.Write([]byte(headder["Sec-WebSocket-Key"]))
		sha.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
		b64 := base64.StdEncoding.EncodeToString(sha.Sum(nil))
		retstr := "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: " + b64 + "\r\n\r\n"
		conn.Write([]byte(retstr))
		fmt.Println("hand success")

		//Start proxy data
		tcpConn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("connect 8080 fail.")
			return
		}
		fmt.Println("proxy connect success!")
		buff_reader := bufio.NewReader(conn)
		num_w, err := buff_reader.WriteTo(tcpConn)
		if err != nil {
			fmt.Println("proxy write err.", num_w)
			tcpConn.Close()
		}

	} else {
		fmt.Println("WebSocket Connection hand fail！")
		return
	}
}

//Tcp To WebSocket
func onTcpAccept(conn net.Conn, addr string) {
	wsConn, err := net.Dial("tcp", addr)
	defer wsConn.Close()
	if err != nil {
		fmt.Println("websocket server connect err.")
		return
	}
	fmt.Println("start proxy connect.")
	header_str := "GET ws://" + addr + " HTTP/1.1\r\nHost: " + addr + "\r\nConnection: Upgrade\r\nPragma: no-cache\r\nCache-Control: no-cache\r\nUpgrade: websocket\r\nOrigin: http://" + addr + "\r\nSec-WebSocket-Version: 13\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.90 Safari/537.36\r\nDNT: 1\r\nAccept-Encoding: gzip, deflate, sdch, br\r\nAccept-Language: zh-CN,zh;q=0.8\r\nCookie: Phpstorm-30b25437=dd95a55a-ff9d-446f-8c6c-731561023f77; wp-settings-time-1=1479280816\r\nSec-WebSocket-Key: ZyK2MTQhx0PozU0K8xy5pA==\r\nSec-WebSocket-Extensions: permessage-deflate; client_max_window_bits\r\n\r\n"
	wsConn.Write([]byte(header_str))
	buff_reader := bufio.NewReader(conn)
	num_w, err := buff_reader.WriteTo(wsConn)
	if err != nil {
		fmt.Println("proxy write err.", num_w)
		wsConn.Close()
		conn.Close()
	}
}
