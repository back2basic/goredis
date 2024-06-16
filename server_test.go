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

	"github.com/redis/go-redis/v9"
	"github.com/tidwall/resp"
)

func TestRedisTTL(t *testing.T) {
	listenaddr := ":5001"
	server := NewServer(Config{
		ListenAddr: listenaddr,
	})
	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(1 * time.Second)


	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost%s", ":5001"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	key := "foo"
	val := "bar"

	if err := rdb.Set(context.Background(), key, val, 10*time.Second).Err(); err != nil {
		t.Fatal(err)
	}

	newVal, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Found key %s with value %s\n", key, newVal)
	fmt.Println("Sleeping for 10 seconds")
	time.Sleep(10 * time.Second)

	newVal2, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		t.Fatal(err)
	}
	if newVal2 == newVal {
		t.Fatalf("expected nil, got %s", newVal2)

	}
}

func TestOfficialRedisClient(t *testing.T) {
	listenaddr := ":5001"
	server := NewServer(Config{
		ListenAddr: listenaddr,
	})
	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(1 * time.Second)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost%s", ":5001"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// key := "foo"
	// val := "bar"

	testCases := map[string]string{
		"foo":  "bar",
		"r":    "vh",
		"your": "mom",
		"step": "dad",
	}

	for key, val := range testCases {
		if err := rdb.Set(context.Background(), key, val, time.Minute).Err(); err != nil {
			t.Fatal(err)
		}

		newVal, err := rdb.Get(context.Background(), key).Result()
		if err != nil {
			t.Fatal(err)
		}

		if newVal != val {
			t.Fatalf("expected %s, got %s", val, newVal)
		}

	}
}

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
