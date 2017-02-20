package server

type ParseJSONFailed struct{}

func (*ParseJSONFailed) Code() string {
	return "parse_json_failed"
}
