package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func initClient(name string, in chan message, out chan message) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8282", Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln("dial error:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)

		for {
			msg := message{}
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Fatalln("recv error:", err)
				return
			}
			log.Printf("[%s] recv: %v\n", name, msg)

			out <- msg
		}
	}()

	for {
		select {
		case msg := <-in:
			err := c.WriteJSON(&msg)
			if err != nil {
				log.Fatalln("send error:", err)
			}
			log.Printf("---------------\n[%s] send: %v\n", name, msg)
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("close client:", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
		}
	}
}
