package account

import (
	"encoding/json"
	"fmt"

	"strconv"

	"github.com/go-resty/resty"
)

func (state *accountActor) session(gameID int) (*SessionSuccess, error) {
	resp, err := resty.R().
		SetFormData(map[string]string{
			"auth_submit":   state.typeAPI + "_CLIENT_SESSION",
			"auth_username": state.account,
			"auth_password": state.password,
			"game_id":       strconv.Itoa(gameID),
		}).
		Post(state.urlAPI)
	if err != nil {
		return nil, fmt.Errorf("Request fail: %s", err)
	}

	res := &struct {
		Result      string `json:"result"`
		Code        int    `json:"code"`
		SessionUUID string `json:"session_uuid"`
		GameUUID    string `json:"game_uuid"`
		Host        string `json:"host"`
	}{}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		return nil, fmt.Errorf("Unmarshal fail: %s", err)
	}

	if res.Result == "Error" {
		return nil, fmt.Errorf("Error: %d", res.Code)
	}

	return &SessionSuccess{res.SessionUUID, res.GameUUID, res.Host}, nil
}
