package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	localAddr := "127.0.0.1:8085"
	remoteAddr := "127.0.0.1:8081"
	local, err := net.Listen("tcp", localAddr)
	if local == nil {
		fatal("cannot listen: %v", err)
	}
	for {
		conn, err := local.Accept()
		if conn == nil {
			fatal("accept failed: %v", err)
		}
		remote, err := net.Dial("tcp", remoteAddr)
		if remote == nil {
			fmt.Fprintf(os.Stderr, "remote dial failed: %v\n", err)
			return
		}
		go forward(conn, remote)
	}
}

func forward(local net.Conn, remote net.Conn) {
	for {
		message, _ := bufio.NewReader(local).ReadString('\n')
		fmt.Print("Proxy Message:", string(message))
		remote.Write([]byte(message + "\n"))

		response, _ := bufio.NewReader(remote).ReadString('\n')
		fmt.Print("Server Response:", string(response))
		local.Write([]byte("PROXIED: " + response + "\n"))
	}
}

func fatal(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "netfwd: %s\n", fmt.Sprintf(s, a))
	os.Exit(2)
}
