package main

import (
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.On("request", func(msg string) {
			log.Printf("Receive %s\n", msg)
			log.Println("emit:", so.Emit("response", "pong"))
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/", server)
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
