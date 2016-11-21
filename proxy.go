package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	argv := os.Args
	if len(argv) < 4 {
		fmt.Println("Argument err, use port proxy example:")
		fmt.Println("root@host:~$./proxy.exe tcp 127.0.0.1:80 192.168.1.102:8080")
		fmt.Println("Source from https://www.github.com/JamesWone/go-port-proxy")
		return
	}
	listen, err := net.Listen(argv[1], argv[2])
	if err != nil {
		fmt.Println("Port listen err:", err.Error())
		listen.Close()
		return
	}
	for {
		client, err := listen.Accept()
		if err != nil {
			fmt.Println("Client->Proxy Connect Fail:", client.RemoteAddr())
			client.Close()
			break
		}

		server, err := net.Dial(argv[1], argv[3])
		if err != nil {
			fmt.Println("Proxy->Server Connect Fail:", err.Error())
			server.Close()
			client.Close()
			break
		}

		fmt.Println("New connection:", client.RemoteAddr().Network(), client.RemoteAddr(), "<->", argv[2], "<->", server.RemoteAddr())
		go proxy(client, server)
		go proxy(server, client)
	}
}

func proxy(client net.Conn, server net.Conn) {
	reader := bufio.NewReader(client)
	_, err := reader.WriteTo(server)
	if err != nil {
		fmt.Println("Proxy err,", err.Error())
		client.Close()
		server.Close()
	}
}
