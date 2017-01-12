package account

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
)

type Game struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Thumb string `json:"thumb"`
}

type Category struct {
	Title string `json:"title"`
	Games []Game `json:"games"`
}

func auth(account, password string) (*AuthSuccess, error) {
	resp, err := resty.R().
		SetFormData(map[string]string{
			"auth_submit":   config.AccountAPI().Type + "_CLIENT_AUTH",
			"auth_username": account,
			"auth_password": password,
		}).
		Post(config.AccountAPI().URL)
	if err != nil {
		return nil, fmt.Errorf("Request fail: %s", err)
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
		return nil, fmt.Errorf("Unmarshal fail: %s", err)
	}

	if res.Result == "Error" {
		return nil, fmt.Errorf("Error: %d", res.Code)
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

	return &AuthSuccess{categories}, nil
}
