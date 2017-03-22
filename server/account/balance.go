package account

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
)

func getBalance(username, password string) Message {
	conf := config.AccountAPI()

	resp, err := resty.R().
		SetFormData(map[string]string{
			"auth_submit":   conf.Type + "_CLIENT_BALANCE",
			"auth_username": username,
			"auth_password": password,
		}).
		Post(conf.URL)
	if err != nil {
		log.Error("get balance failed",
			zap.Error(err),
		)
		return &GetBalanceFailed{}
	}

	res := &struct {
		Result  string  `json:"result"`
		Code    int     `json:"code"`
		Balance float64 `json:"balance"`
	}{}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		log.Error("get balance failed",
			zap.Error(err),
		)
		return &GetBalanceFailed{}
	}

	if res.Result == "Error" {
		log.Error("get balance failed",
			zap.String("result", res.Result),
			zap.Int("code", res.Code),
		)
		return &GetBalanceFailed{}
	}

	return &Balance{res.Balance}
}
