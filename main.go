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

	time.Sleep(1 * time.Second)

	appIn := make(chan message)
	appOut := make(chan message)

	go initClient("app", appIn, appOut)

	webIn := make(chan message)
	webOut := make(chan message)

	go initClient("web", webIn, webOut)

	appIn <- message{
		Type: "command.subscribe",
		Data: messages.Subscribe{
			Topics: []string{
				"event.login.success",
			},
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
						Topics: []string{
							"event.login.success",
							"event.account.fail",
							"event.account.auth.success",
							"event.account.auth.fail",
							"event.account.balance.success",
							"event.account.balance.fail",
							"event.account.session.success",
							"event.account.session.fail",
							"event.account.withdraw.success",
							"event.account.withdraw.fail",
						},
					},
				}
			}
		case msg := <-webOut:
			switch msg.Type {
			case "event.subscribe.success":
				webIn <- message{
					Type: "command.login",
					Data: session.Login{client},
				}
			case "event.login.success":
				webIn <- message{
					Type: "command.account.auth",
					Data: map[string]interface{}{
						"account":  "1191100006",
						"password": "3129",
					},
				}
			case "event.account.auth.success":
				webIn <- message{
					Type: "command.account.balance",
				}
			case "event.account.balance.success":
				webIn <- message{
					Type: "command.account.session",
					Data: map[string]interface{}{
						"game_id": 83,
					},
				}
			case "event.account.session.success":
				webIn <- message{
					Type: "command.account.withdraw",
				}
			}
		}
	}

	//console.ReadLine()
}
