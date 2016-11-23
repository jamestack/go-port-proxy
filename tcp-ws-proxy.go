package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var ws = websocket.Upgrader{}

var laddr string = "localhost:8087"
var raddr string = "localhost:8080"

func main() {
	args := os.Args
	if len(args) < 4 {
		fmt.Println("tcp proxy websocket")
		fmt.Println("root@host:./proxy.exe tcp-to-ws localhost:8087 www.jamestack.tk:8080")
		return
	}
	laddr, raddr = args[2], args[3]
	switch args[1] {
	case "ws-to-tcp":
		ws_to_tcp()
	case "tcp-to-ws":
		tcp_to_ws()
	default:
		fmt.Println("参数错误!")
	}
}

func tcp_to_ws() {
	listen, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Println("port listen err:", err)
		listen.Close()
		return
	}
	for {
		client, err := listen.Accept()
		if err != nil {
			log.Println("tcp connect err:", client.RemoteAddr())
			client.Close()
			break
		}

		go proxy_tcp_to_ws(client)
	}
}

func proxy_tcp_to_ws(client net.Conn) {
	defer client.Close()

	ws, _, err := websocket.DefaultDialer.Dial("ws://"+raddr+"/proxy", nil)
	defer ws.Close()
	if err != nil {
		log.Println("dial err:", err)
		return
	}

	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("ws read err:", err)
				break
			}
			//			log.Println("ws read:", message)
			_, err = client.Write(message)
			if err != nil {
				log.Println("tcp write err:", err)
				break
			}
		}
	}()

	for {
		buff := make([]byte, 65535)
		num, err := client.Read(buff)
		if err != nil {
			log.Println("tcp read err:", err)
			break
		}
		err = ws.WriteMessage(websocket.BinaryMessage, buff[:num])
		if err != nil {
			log.Println("ws wrait err:", err)
			break
		}
	}
}

func ws_to_tcp() {
	http.HandleFunc("/", handle_root)
	http.HandleFunc("/proxy", handle_ws)
	err := http.ListenAndServe(laddr, nil)
	if err != nil {
		log.Println("http listen err:", err.Error())
		return
	}
}

func handle_root(w http.ResponseWriter, r *http.Request) {
	log.Println("/")
	html := "Hello Websocket Proxy!" +
		"<br/>" +
		"Source by http://github.com/jameswone/websocket-proxy"
	_, err := w.Write([]byte(html))
	if err != nil {
		log.Println("html write err:", err)
	}
}

func handle_ws(w http.ResponseWriter, r *http.Request) {
	log.Println("/proxy")
	ws, err := ws.Upgrade(w, r, nil)
	if err != nil {
		log.Print("ws upgrade err:", err)
		return
	}
	defer ws.Close()

	client, err := net.Dial("tcp", raddr)
	if err != nil {
		log.Println("tcp dial err:", err)
		return
	}
	defer client.Close()

	go func() {
		for {
			buff := make([]byte, 65535)
			num, err := client.Read(buff)
			if err != nil {
				log.Println("tcp read err:", err)
				break
			}
			err = ws.WriteMessage(websocket.BinaryMessage, buff[:num])
			if err != nil {
				log.Println("ws wrait err:", err)
				break
			}
		}
	}()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("ws read err:", err)
			break
		}
		//		log.Println("ws read:", message)
		_, err = client.Write(message)
		if err != nil {
			log.Println("tcp write err:", err)
			break
		}
	}
}
