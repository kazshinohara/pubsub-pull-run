package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Message struct {
	Greeting string `json:"greeting"`
}

var (
	port    = os.Getenv("PORT")
	projectId = os.Getenv("PROJECT_ID")
	topicId   = os.Getenv("TOPIC_ID")
)

func publishMsgs(number int) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		fmt.Print("client error.err:", err)
		os.Exit(1)
	}

	event := &Message{
		Greeting: "Hello, Cloud Run",
	}

	msg, _ := json.Marshal(event)
	t := client.Topic(topicId)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})

	id, err := result.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Published No.%d message; msg ID: %v\n", number, id)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I'm healthy.")
}

func main() {
	go func(){
		for i := 0; ; i++ {
			time.Sleep(1 * time.Second)
			publishMsgs(i)
		}
	}()
	http.HandleFunc("/", handler)
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
