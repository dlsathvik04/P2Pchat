package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go listener(&wg)
	time.Sleep(time.Second)
	go dialer(&wg)
	wg.Wait()

}

func dialer(wg *sync.WaitGroup) {
	peer := Peer{
		StunUrl:          "http://localhost:8000",
		Port:             9000,
		Username:         "dls",
		password:         "124j",
		handleConnection: handleConnection,
	}
	conn, err := peer.Dial("dlsathvik04")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.Write([]byte("Hello"))
	wg.Done()
}

func listener(wg *sync.WaitGroup) {
	peer := Peer{
		StunUrl:          "http://localhost:8000",
		Port:             8500,
		Username:         "dlsathvik04",
		password:         "1234",
		handleConnection: handleConnection,
	}
	peer.StartAccepting()
	wg.Done()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	fmt.Println("Received data:", string(buf[:n]))
	_, err = conn.Write([]byte("Hello from server"))
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		return
	}

	fmt.Println("Response sent to client")
}
