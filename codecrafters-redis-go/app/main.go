package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var storage = map[string]string{}
var expired = map[string]time.Time{}
var list = map[string][]string{}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			os.Exit(1)
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	for {
		first, err := r.ReadString('\n')
		if err != nil {
			return
		}

		fmt.Println("first value should be array length", first)

		var length int
		fmt.Sscanf(first, "*%d", &length)

		msg := make([]string, 0)
		for _ = range length {
			l, err := r.ReadString('\n')
			if err != nil {
				return
			}
			var ln int
			fmt.Sscanf(l, "$%d", &ln)

			cmd := make([]byte, ln)
			n, err := r.Read(cmd)
			if err != nil {
				return
			}

			msg = append(msg, string(cmd[:n]))

			r.ReadString('\n')
		}

		fmt.Println("Message: ", strings.Join(msg, " "))

		cmd := strings.ToUpper(msg[0])
		switch cmd {
		case "ECHO":
			res := fmt.Sprintf("$%d\r\n%s\r\n", len(msg[1]), msg[1])
			conn.Write([]byte(res))
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		case "SET":
			key := msg[1]
			value := msg[2]
			if len(msg) > 3 {
				var expTimeInt int
				fmt.Sscanf(msg[4], "%d", &expTimeInt)
				expired[key] = time.Now().Add(time.Duration(expTimeInt) * time.Millisecond)
			}

			storage[key] = value

			conn.Write([]byte("+OK\r\n"))
		case "GET":
			key := msg[1]
			val := storage[key]
			t, ok := expired[key]
			if ok && time.Now().After(t) {
				conn.Write([]byte("$-1\r\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
			}
		case "RPUSH":
			key := msg[1]
			values := make([]string, 0)
			for i := 2; i < len(msg); i++ {
				values = append(values, msg[i])
			}

			l, ok := list[key]
			if !ok {
				list[key] = values
			} else {
				list[key] = append(list[key], values...)
			}

			conn.Write([]byte(fmt.Sprintf(":%d\r\n", len(l)+len(values))))
		case "LRANGE":

		}
	}
}
