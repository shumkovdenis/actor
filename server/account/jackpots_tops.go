package account

import (
	"encoding/json"

	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
	"github.com/shumkovdenis/club/logger"
	"go.uber.org/zap"
)

func jackpotsTops(username, password string) Message {
	conf := config.AccountAPI()

	resp, err := resty.R().
		SetFormData(map[string]string{
			"auth_submit":   conf.Type + "_CLIENT_JACKPOTS",
			"auth_username": username,
			"auth_password": password,
		}).
		Post(conf.URL)
	if err != nil {
		logger.L().Error("jackpots tops failed",
			zap.Error(err),
		)
		return &GetJackpotsTopsFailed{}
	}

	res := &struct {
		Result   string `json:"result"`
		Code     int    `json:"code"`
		Jackpots []struct {
			Account string `json:"account"`
			Win     string `json:"win"`
			Moment  string `json:"moment"`
		} `json:"jackpots"`
	}{}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		logger.L().Error("jackpots tops failed",
			zap.Error(err),
		)
		return &GetJackpotsTopsFailed{}
	}

	if res.Result == "Error" {
		logger.L().Error("jackpots tops failed",
			zap.String("result", res.Result),
			zap.Int("code", res.Code),
		)
		return &GetJackpotsTopsFailed{}
	}

	jackpots := make([]Jackpot, len(res.Jackpots))
	for i, jackpot := range res.Jackpots {
		jackpots[i] = Jackpot{
			Account: jackpot.Account,
			Win:     jackpot.Win,
			Date:    jackpot.Moment,
		}
	}

	return &JackpotsTops{jackpots}
}
