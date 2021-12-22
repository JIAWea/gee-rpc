package main

import (
	"context"
	"gee-rpc"
	"log"
	"net"
	"sync"
	"time"
)

type Foo int

type Args struct {
	Num1, Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func startServer(addr chan<- string) {
	var foo Foo
	if err := geerpc.Register(&foo); err != nil {
		log.Fatal("register error: ", err)
	}

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
			// args := fmt.Sprintf("geerpc req %d", i)
			args := &Args{Num1: i, Num2: i * i}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			var reply int
			if err := client.Call(ctx, "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
			cancel()
		}(i)
	}
	wg.Wait()
}
