package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/eiannone/keyboard"
)

var (
	conn_active bool = true
)

func main() {
	host := flag.String("h", "localhost", "Host")
	port := flag.Int("p", 12321, "Port")
	flag.Parse()
	startClient(fmt.Sprintf("%s:%d", *host, *port))
}

func processClient(conn net.Conn) {
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		fmt.Println(err)
	}
	conn.Close()
	conn_active = false
	fmt.Printf("connection closed. press key to quit...\n")
}

func startClient(addr string) {
	fmt.Printf("connecting to: %s ...\n", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("can't connect to server: %s\n", err)
		return
	}

	go processClient(conn)

	kerr := keyboard.Open()
	if kerr != nil {
		panic(kerr)
	}
	defer keyboard.Close()

	fmt.Println("press ESC to quit")
	var out []byte;
	for conn_active {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		} else if key == keyboard.KeyEsc {
			break
		} else if key == keyboard.KeyEnter {
			out = []byte{0x0d, 0x0a}
		} else if key == keyboard.KeySpace {
			out = []byte{0x20}
		} else if key == keyboard.KeyBackspace {
			out = []byte{0x08}
		} else {
			out = []byte{byte(char)}
		}
		conn.Write(out)
	}
}

