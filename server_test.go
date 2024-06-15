package main

import (
	"bytes"
	"context"
	"fmt"
	"goredis/client"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/tidwall/resp"
)

func TestHelloServer(t *testing.T) {
	buf := &bytes.Buffer{}
	rw := resp.NewWriter(buf)
	rw.WriteString("OK")
	fmt.Println(buf.String())

	in := map[string]string{
		"first":  "1",
		"second": "2",
	}
	out := respWriteMap(in)
	fmt.Println(out)
}

func TestServerWithMultipleClients(t *testing.T) {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(1 * time.Second)

	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			c, err := client.New(("localhost:5001"))
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()

			key := fmt.Sprintf("client_foo_%d", it)
			value := fmt.Sprintf("client_bar_%d", it)
			if err := c.Set(context.TODO(), key, value); err != nil {
				log.Fatal(err)
			}

			val, err := c.Get(context.TODO(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Client %d GET => %s\n", it, val)
			wg.Done()
		}(i)
	}

	wg.Wait()

	time.Sleep(1 * time.Second)
	if len(server.peers) != 0 {
		t.Fatalf("expected 0 but got %d", len(server.peers))
	}

}
