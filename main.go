package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

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

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("cound not open %s: %s\n", inputFilePath, err)
	}

	strCh := getLinesChannel(file)

	for v := range strCh {
		fmt.Println("read:", v)
	}

}
