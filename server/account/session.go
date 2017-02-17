package account

import (
	"encoding/json"
	"strconv"

	"go.uber.org/zap"

	"github.com/go-resty/resty"
	"github.com/shumkovdenis/club/config"
)

func getGameSession(username, password string, gameID int) Message {
	conf := config.AccountAPI()

	resp, err := resty.R().
		SetFormData(map[string]string{
			"auth_submit":   conf.Type + "_CLIENT_SESSION",
			"auth_username": username,
			"auth_password": password,
			"game_id":       strconv.Itoa(gameID),
		}).
		Post(conf.URL)
	if err != nil {
		log.Error("get session failed",
			zap.Error(err),
		)
		return &GetGameSessionFailed{}
	}

	res := &struct {
		Result      string `json:"result"`
		Code        int    `json:"code"`
		SessionUUID string `json:"session_uuid"`
		GameUUID    string `json:"game_uuid"`
		Host        string `json:"host"`
	}{}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		log.Error("get session failed",
			zap.Error(err),
		)
		return &GetGameSessionFailed{}
	}

	if res.Result == "Error" {
		log.Error("get session failed",
			zap.String("result", res.Result),
		)
		return &GetGameSessionFailed{}
	}

	return &GameSession{
		SessionID: res.SessionUUID,
		GameID:    res.GameUUID,
		ServerURL: res.Host,
	}
}
