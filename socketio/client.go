package main

import (
	"log"

	"github.com/zhouhui8915/go-socket.io-client"
)

func main() {

	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	uri := "http://127.0.0.1:8000/"

	client, err := socketio_client.NewClient(uri, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}

	client.On("response", func(msg string) {
		log.Println("emit:", client.Emit("request", "ping"))
	})

	client.Emit("request", "ping")

	for {

	}
}
