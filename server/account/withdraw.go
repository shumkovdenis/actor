package account

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
)

func withdraw(username, password string) Message {
	conf := config.AccountAPI()

	resp, err := resty.R().
		SetFormData(map[string]string{
			"auth_submit":   conf.Type + "_CLIENT_WITHDRAW",
			"auth_username": username,
			"auth_password": password,
		}).
		Post(conf.URL)
	if err != nil {
		logger.L().Error("withdraw failed",
			zap.Error(err),
		)
		return &WithdrawFailed{}
	}

	res := &struct {
		Result string `json:"result"`
		Code   int    `json:"code"`
	}{}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		logger.L().Error("withdraw failed",
			zap.Error(err),
		)
		return &WithdrawFailed{}
	}

	if res.Result == "Error" {
		logger.L().Error("withdraw failed",
			zap.String("result", res.Result),
			zap.Int("code", res.Code),
		)
		return &WithdrawFailed{}
	}

	return &WithdrawSuccess{}
}
