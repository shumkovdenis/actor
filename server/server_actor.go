package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/labstack/echo"
	"github.com/shumkovdenis/club/config"
)

type renderer struct {
	template *template.Template
}

func (r *renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.template.ExecuteTemplate(w, name, data)
}

type serverActor struct{}

func newServerActor() actor.Actor {
	return &serverActor{}
}

func (state *serverActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
		conf := config.Server()
		publicPath := conf.PublicPath

		e := echo.New()
		e.Renderer = &renderer{
			template: template.Must(template.ParseGlob(publicPath + "/*.html")),
		}
		e.Static("/", publicPath)
		e.GET("/", index)

		props := actor.FromInstance(newAPIActor(e.Group("/api")))
		ctx.SpawnNamed(props, "api")

		props = actor.FromInstance(newConnManagerActor(e.Group("/conn")))
		ctx.SpawnNamed(props, "conns")

		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Port)))
	}
}

func index(c echo.Context) error {
	conf := config.Server()

	session := c.QueryParam("session")

	state := &struct {
		Websocket string
		Session   string
	}{
		Websocket: conf.WebSocketURL(),
		Session:   session,
	}

	return c.Render(http.StatusOK, "index", state)
}
