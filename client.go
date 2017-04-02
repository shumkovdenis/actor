package club

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
	"github.com/shumkovdenis/club/logger"

	"io/ioutil"

	"fmt"
)

type message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type data struct {
	States map[string]message `json:"states"`
	Events map[string]string  `json:"events"`
}

func StartClient(dataFile string) error {
	if len(strings.TrimSpace(dataFile)) == 0 {
		return errors.New("must specify path to data file")
	}

	bytes, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return err
	}

	d := &data{}

	if err := json.Unmarshal(bytes, d); err != nil {
		return err
	}

	if err := connect(d); err != nil {
		return err
	}

	return nil
}

func connect(d *data) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8282", Path: "/conn/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	done := make(chan struct{})
	in := make(chan message, 1)

	go func() {
		defer c.Close()
		defer close(done)

		for {
			msg := message{}
			err := c.ReadJSON(&msg)
			if err != nil {
				logger.L().Fatal("recv error:", zap.Error(err))
			}
			logger.L().Info(fmt.Sprintf("recv: %v\n", msg))

			if nextState, ok := d.Events[msg.Type]; ok {
				if nextCommand, ok := d.States[nextState]; ok {
					in <- nextCommand
				}
			}
		}
	}()

	in <- d.States["#start"]

	for {
		select {
		case msg := <-in:
			err := c.WriteJSON(&msg)
			if err != nil {
				logger.L().Fatal("send error:", zap.Error(err))
			}
			logger.L().Info(fmt.Sprintf("\n---------------\nsend: %v\n", msg))
		case <-interrupt:
			logger.L().Info("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.L().Info("close client:", zap.Error(err))
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
		}
	}

	return nil
}

/*func initClient(name string, in chan message, out chan message) {
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
*/
