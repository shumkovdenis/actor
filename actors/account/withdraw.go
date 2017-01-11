package account

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty"
)

func (state *accountActor) withdraw() (*WithdrawSuccess, error) {
	resp, err := resty.R().
		SetFormData(map[string]string{
			"auth_submit":   state.typeAPI + "_CLIENT_WITHDRAW",
			"auth_username": state.account,
			"auth_password": state.password,
		}).
		Post(state.urlAPI)
	if err != nil {
		return nil, fmt.Errorf("Request fail: %s", err)
	}

	res := &struct {
		Result string `json:"result"`
		Code   int    `json:"code"`
	}{}
	if err = json.Unmarshal(resp.Body(), res); err != nil {
		return nil, fmt.Errorf("Unmarshal fail: %s", err)
	}

	if res.Result == "Error" {
		return nil, fmt.Errorf("Error: %d", res.Code)
	}

	return &WithdrawSuccess{}, nil
}
