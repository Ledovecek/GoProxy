package main

import (
	"net"
)

const (
	PROXY_HOST  = "localhost"
	PROXY_PORT  = "90"
	SERVER_HOST = "localhost"
	SERVER_PORT = "25565"
)

func main() {
	listener, _ := net.Listen("tcp", PROXY_HOST+":"+PROXY_PORT)
	print("[PROXY] Running on port: ", PROXY_PORT)
	handleCon(listener)
}

func handleCon(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	dial, err := net.Dial("tcp", SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		return
	}
	buffer := make([]byte, 256)
	go handleIncomingClientRequest(&dial, &conn)
	for {
		bytesCount, err := dial.Read(buffer)
		if err != nil {
			return
		}
		count, err := conn.Write(buffer[0:bytesCount])
		if count != bytesCount || err != nil {
			conn.Close()
			dial.Close()
			return
		}
	}
}

func handleIncomingClientRequest(serverConnection *net.Conn, clientConnection *net.Conn) {
	buffer := make([]byte, 256)
	for {
		bytesCount, err := (*clientConnection).Read(buffer)
		if err != nil {
			return
		}
		count, err := (*serverConnection).Write(buffer[0:bytesCount])
		if count != bytesCount || err != nil {
			(*serverConnection).Close()
			(*clientConnection).Close()
			return
		}
	}
}
