package server

import "github.com/shumkovdenis/club/server/core"

type ParseJSONFailed struct{}

func (*ParseJSONFailed) Code() string { return "parse_json_failed" }

func failResp(fail core.Fail) interface{} {
	return &struct {
		Code string `json:"code"`
	}{
		Code: fail.Code(),
	}
}
