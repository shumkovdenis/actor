package account

import (
	"encoding/json"

	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
	"go.uber.org/zap"
)

func authorize(account, password string) Outgoing {
	conf := config.AccountAPI()

	resp, err := resty.R().
		SetFormData(map[string]string{
			"auth_submit":   conf.Type + "_CLIENT_AUTH",
			"auth_username": account,
			"auth_password": password,
		}).
		Post(conf.URL)
	if err != nil {
		log.Error("authorization failed",
			zap.Error(err),
		)
		return &AuthorizationFailed{}
	}

	res := &struct {
		Result string `json:"result"`
		Code   int    `json:"code"`
		Groups []struct {
			Title string `json:"title"`
			Games []struct {
				ID    string `json:"id"`
				Title string `json:"title"`
			}
		} `json:"groups"`
	}{}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		log.Error("authorization failed",
			zap.Error(err),
		)
		return &AuthorizationFailed{}
	}

	if res.Result == "Error" {
		log.Error("authorization failed",
			zap.String("result", res.Result),
		)
		return &AuthorizationFailed{}
	}

	categories := make([]Category, len(res.Groups))
	for i, group := range res.Groups {
		games := make([]Game, len(group.Games))
		for j, game := range group.Games {
			games[j] = Game{
				ID:    game.ID,
				Title: game.Title,
			}
		}
		categories[i] = Category{
			Title: group.Title,
			Games: games,
		}
	}

	return &Authorized{categories}
}
