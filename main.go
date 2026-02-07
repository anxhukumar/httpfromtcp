package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":42069"

func main() {
	tcpListener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer tcpListener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("A tcp connection has been accepted...")

		msgCh := getLinesChannel(conn)
		for v := range msgCh {
			fmt.Printf("%s\n", v)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {

	strCh := make(chan string)

	go func() {
		defer f.Close()
		defer close(strCh)

		currentLine := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				break
			}
			parts := strings.Split(string(buffer[:n]), "\n")
			for i, line := range parts {
				if i == len(parts)-1 {
					currentLine += line
					break
				}
				strCh <- fmt.Sprintf("%s%s", currentLine, line)
				currentLine = ""
			}
		}

		if len(currentLine) > 0 {
			strCh <- currentLine
		}
	}()
	return strCh
}
