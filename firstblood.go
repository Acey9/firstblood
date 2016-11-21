package main

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	args := os.Args[1:]
	if args == nil || len(args) < 2 {
		log("arguments error...")
		return
	}

	listenIP := args[0]
	for _, port := range args[1:] {
		addr := listenIP + ":" + port
		go listen("tcp", addr)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

forever:
	for {
		select {
		case <-sig:
			fmt.Println("Interrupt signal recevied, stopping")
			break forever
		}
	}
}

func listen(network, address string) error {
	tel, err := net.Listen(network, address)
	if err != nil {
		log(err)
		return err
	}

	for {
		conn, err := tel.Accept()
		if err != nil {
			log(err)
			break
		}
		go initHandler(conn)
	}
	log("Stopped accepting:", address)
	return nil
}

func initHandler(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	buf := make([]byte, 2048)
	l, err := conn.Read(buf)
	if err != nil || l <= 0 {
		return
	}
	encoded := base64.StdEncoding.EncodeToString(buf[:l])
	log(conn.RemoteAddr().String(), conn.LocalAddr().String(), encoded)
}

func log(a ...interface{}) {
	fmt.Printf("[%s] ", time.Now().String())
	fmt.Println(a...)
}
