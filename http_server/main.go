package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mdlayher/vsock"
)

func main() {
	const port uint32 = 8000

	ln, err := vsock.Listen(port, nil)
	if err != nil {
		log.Fatalf("vsock listen: %v", err)
	}
	log.Printf("enclave HTTP server listening on vsock::%d", port)

	for {
		c, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go func() {
			defer c.Close()
			c.SetDeadline(time.Now().Add(30 * time.Second))
			br := bufio.NewReader(c)
			reqLine, _ := br.ReadString('\n')
			// best-effort read headers (not strictly needed)
			for {
				h, _ := br.ReadString('\n')
				if strings.TrimSpace(h) == "" {
					break
				}
			}
			body := "Hello from inside the Nitro Enclave!\nYou reached " + strings.TrimSpace(reqLine) + "\n"
			fmt.Fprintf(c, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
		}()
	}
}
