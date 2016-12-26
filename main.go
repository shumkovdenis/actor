package main

import (
	"log"

	"github.com/shumkovdenis/actor/actors/server"
	"github.com/shumkovdenis/actor/actors/session"
	"github.com/shumkovdenis/actor/messages"

	"time"

	"github.com/AsynkronIT/gam/actor"
)

func main() {
	log.SetFlags(0)

	props := actor.FromProducer(server.New)
	go actor.Spawn(props)

	time.Sleep(2 * time.Second)

	appIn := make(chan message)
	appOut := make(chan message)

	go initClient("app", appIn, appOut)

	webIn := make(chan message)
	webOut := make(chan message)

	go initClient("web", webIn, webOut)

	appIn <- message{
		Type: "command.subscribe",
		Data: messages.Subscribe{
			Topic: "event.login.success",
		},
	}

	var client string

	for {
		select {
		case msg := <-appOut:
			switch msg.Type {
			case "event.subscribe.success":
				appIn <- message{
					Type: "command.login",
					Data: session.Login{},
				}
			case "event.login.success":
				data := msg.Data.(map[string]interface{})
				client = data["client"].(string)
				webIn <- message{
					Type: "command.subscribe",
					Data: messages.Subscribe{
						Topic: "event.login.success",
					},
				}
			}
		case msg := <-webOut:
			if msg.Type == "event.subscribe.success" {
				webIn <- message{
					Type: "command.login",
					Data: session.Login{client},
				}
			}
		}
	}

	//console.ReadLine()
}
