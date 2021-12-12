package main

import (
	"encoding/json"
	"fmt"
	"gee-rpc"
	"gee-rpc/codec"
	"log"
	"net"
	"time"
)

func startServer(addr chan<- string) {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("network error: %v", err)
	}
	log.Println("start rpc server on ", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func  main() {
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple client
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	// send options
	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	cc := codec.NewGobCodec(conn)
	// send request and receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply: ", reply)
	}
}
