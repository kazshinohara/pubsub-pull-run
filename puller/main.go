package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	port    = os.Getenv("PORT")
	projectId = os.Getenv("PROJECT_ID")
	subId     = os.Getenv("SUB_ID")
)

type Message struct {
	Greeting string `json:"greeting"`
}

func pullMsgs(projectID, subID string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	var mu sync.Mutex
	sub := client.Subscription(subID)
	cctx, _ := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		var receivedMessage Message
		err := json.Unmarshal(msg.Data, &receivedMessage)
		if err != nil {
			log.Printf("json unmarshal: %v", err)
		}
		log.Printf("Greeting message: %q\n", receivedMessage.Greeting)
		msg.Ack()
	})
	if err != nil {
		log.Printf("Receive: %v", err)
	}
	return nil
}

func fib(n int) int {
	if n == 0 || n == 1 {
		return n
	}
	return fib(n-2) + fib(n-1)
}

func fibHandler(w http.ResponseWriter, r*http.Request) {
	go func() {
		for i := 0; i < 1000; i++ {
			log.Printf("fib(%d) = %d\n", i, fib(i))
		}
	}()
	fmt.Fprintf(w, "Now I'm doing fibonacci.")
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I'm healthy.")
}

func main() {
	go pullMsgs(projectId, subId)
	http.HandleFunc("/", handler)
	http.HandleFunc("/fib", fibHandler)
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
