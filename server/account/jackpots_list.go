package account

import (
	"encoding/json"

	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"go.uber.org/zap"
)

func jackpotsList(username, password string) Message {
	conf := config.AccountAPI()

	resp, err := resty.R().
		SetFormData(map[string]string{
			"auth_submit":   conf.Type + "_CLIENT_JACKPOTCOUNTERS",
			"auth_username": username,
			"auth_password": password,
		}).
		Post(conf.URL)
	if err != nil {
		logger.L().Error("jackpots list failed",
			zap.Error(err),
		)
		return &GetJackpotsListFailed{}
	}

	res := &struct {
		Result   string `json:"result"`
		Code     int    `json:"code"`
		Jackpots struct {
			ExtraJackpot string `json:"ExtraJackpot"`
			SuperJackpot string `json:"SuperJackpot"`
			Jackpot      string `json:"Jackpot"`
		} `json:"jackpots"`
	}{}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		logger.L().Error("jackpots list failed",
			zap.Error(err),
		)
		return &GetJackpotsListFailed{}
	}

	if res.Result == "Error" {
		logger.L().Error("jackpots list failed",
			zap.String("result", res.Result),
			zap.Int("code", res.Code),
		)
		return &GetJackpotsListFailed{}
	}

	return &JackpotsList{
		Large:  res.Jackpots.ExtraJackpot,
		Medium: res.Jackpots.SuperJackpot,
		Small:  res.Jackpots.Jackpot,
	}
}
