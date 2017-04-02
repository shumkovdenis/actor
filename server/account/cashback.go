package account

import (
	"encoding/json"

	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"go.uber.org/zap"
)

func cashback(username, password string) Message {
	conf := config.AccountAPI()

	resp, err := resty.R().
		SetFormData(map[string]string{
			"auth_submit":   conf.Type + "_CLIENT_CASHBACK",
			"auth_username": username,
			"auth_password": password,
		}).
		Post(conf.URL)
	if err != nil {
		logger.L().Error("cashback failed",
			zap.Error(err),
		)
		return &CashbackFailed{}
	}

	res := &struct {
		Result   string  `json:"result"`
		Code     int     `json:"code"`
		Cashback float64 `json:"cashback"`
	}{}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		logger.L().Error("cashback failed",
			zap.Error(err),
		)
		return &CashbackFailed{}
	}

	if res.Result == "Error" {
		logger.L().Error("cashback failed",
			zap.String("result", res.Result),
			zap.Int("code", res.Code),
		)
		return &CashbackFailed{}
	}

	return &CashbackSuccess{res.Cashback}
}
