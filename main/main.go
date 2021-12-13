package main

import (
	"fmt"
	"gee-rpc"
	"log"
	"net"
	"sync"
	"time"
)

func startServer(addr chan<- string) {
	l, err := net.Listen("tcp", ":8889")
	if err != nil {
		log.Fatalf("network error: %v", err)
	}
	log.Println("start rpc server on ", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple client
	// conn, _ := net.Dial("tcp", <-addr)
	client, _ := geerpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send options
	// _ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	// cc := codec.NewGobCodec(conn)
	var wg sync.WaitGroup

	// send request and receive response
	for i := 0; i < 5; i++ {
		wg.Add(1)
		// h := &codec.Header{
		// 	ServiceMethod: "Foo.sum",
		// 	Seq:           uint64(i),
		// }
		// _ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		// _ = cc.ReadHeader(h)
		// var reply string
		// _ = cc.ReadBody(&reply)
		// log.Println("reply: ", reply)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply: ", reply)
		}(i)
	}
	wg.Wait()
}
