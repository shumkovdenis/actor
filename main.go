package main

import (
	"crypto/tls"
	"log"

	"github.com/go-resty/resty"
	"github.com/shumkovdenis/actor/actors/rates"
	"github.com/shumkovdenis/actor/actors/server"

	"time"

	"github.com/AsynkronIT/gam/actor"
	"github.com/shumkovdenis/actor/actors/group"
)

func main() {
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	log.SetFlags(0)

	props := actor.FromProducer(server.New)
	go actor.Spawn(props)

	time.Sleep(200 * time.Millisecond)

	props = actor.FromProducer(group.NewActor)
	go actor.SpawnNamed(props, "/rates")

	time.Sleep(200 * time.Millisecond)

	props = actor.FromInstance(rates.NewActor(actor.NewLocalPID("/rates")))
	go actor.Spawn(props)

	time.Sleep(200 * time.Millisecond)

	appIn := make(chan message)
	appOut := make(chan message)

	go initClient("app", appIn, appOut)

	webIn := make(chan message)
	webOut := make(chan message)

	go initClient("web", webIn, webOut)

	appIn <- message{
		Type: "command.subscribe",
		Data: map[string]interface{}{
			"topics": []string{
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
				}
			case "event.login.success":
				data := msg.Data.(map[string]interface{})
				client = data["client"].(string)
				webIn <- message{
					Type: "command.subscribe",
					Data: map[string]interface{}{
						"topics": []string{
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
							"event.rates.change",
						},
					},
				}
			}
		case msg := <-webOut:
			switch msg.Type {
			case "event.subscribe.success":
				webIn <- message{
					Type: "command.login",
					Data: map[string]interface{}{
						"client": client,
					},
				}
				// case "event.login.success":
				// 	webIn <- message{
				// 		Type: "command.account.auth",
				// 		Data: map[string]interface{}{
				// 			"account":  "1191100006",
				// 			"password": "3129",
				// 		},
				// 	}
				// case "event.account.auth.success":
				// 	webIn <- message{
				// 		Type: "command.account.balance",
				// 	}
				// case "event.account.balance.success":
				// 	webIn <- message{
				// 		Type: "command.account.session",
				// 		Data: map[string]interface{}{
				// 			"game_id": 83,
				// 		},
				// 	}
				// case "event.account.session.success":
				// 	webIn <- message{
				// 		Type: "command.account.withdraw",
				// 	}
			}
		}
	}

	//console.ReadLine()
}
