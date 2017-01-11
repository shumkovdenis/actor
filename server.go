package club

import (
	"crypto/tls"

	"github.com/AsynkronIT/gam/actor"
	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/actors/group"
	"github.com/shumkovdenis/club/actors/rates"
	"github.com/shumkovdenis/club/actors/server"
)

func StartServer() error {
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	props := actor.FromProducer(group.NewActor)
	pid := actor.SpawnNamed(props, "/rates")

	props = actor.FromInstance(rates.New(pid))
	actor.Spawn(props)

	props = actor.FromProducer(server.New)
	actor.Spawn(props)

	return nil
}

// func main() {
// if err := app.ReadInfo(); err != nil {
// 	logger.Fatal(err.Error())
// }

// if err := app.ReadConfig(); err != nil {
// 	logger.Fatal(err.Error())
// }

// logger.Info("Run.")
// }

/*
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

	props = actor.FromProducer(group.NewActor)
	go actor.SpawnNamed(props, "/update")

	time.Sleep(200 * time.Millisecond)

	props = actor.FromInstance(update.NewActor(actor.NewLocalPID("/update")))
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
							"event.login.fail",
							"event.join.success",
							"event.join.fail",
							"event.app.update.no",
							"event.app.update.available",
							"event.app.update.download",
							"event.app.update.ready",
							"event.app.update.install",
							"event.app.update.restart",
							"event.app.update.fail",
							"event.account.fail",
							"event.account.auth.success",
							"event.account.auth.fail",
							"event.account.balance.success",
							"event.account.balance.fail",
							"event.account.session.success",
							"event.account.session.fail",
							"event.account.withdraw.success",
							"event.account.withdraw.fail",
							// "event.rates.change",
							// "event.rates.fail",
						},
					},
				}
			}
		case msg := <-webOut:
			switch msg.Type {
			case "event.subscribe.success":
				webIn <- message{
					Type: "command.join",
					Data: map[string]interface{}{
						"client": client,
					},
				}
				// case "event.join.success":
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
				// case "event.account.withdraw.success":
				// 	time.Sleep(10 * time.Second)
				// 	webIn <- message{
				// 		Type: "command.unsubscribe",
				// 		Data: map[string]interface{}{
				// 			"topics": []string{
				// 				"event.rates.change",
				// 			},
				// 		},
				// 	}
			}
		}
	}

	console.ReadLine()
}
*/
