package client

import (
	"context"
	"fmt"
	"sync"

	"log"
	"testing"
)

func TestNewClients(t *testing.T) {
	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			c, err := New(("localhost:5001"))
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
}

func TestNewClient(t *testing.T) {

	c, err := New(("127.0.0.1:5001"))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 1000; i++ {
		fmt.Println("SET =>", fmt.Sprintf("bar_%d", i))
		if err := c.Set(context.TODO(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}

		val, err := c.Get(context.TODO(), fmt.Sprintf("foo_%d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("GET =>", val)
	}

}
