package format_api

import (
	"reflect"
)

type ApiFormat struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ApiJsonWithError(code int64, msg string, data interface{}) (interface{}, error) {
	if data == nil || !reflect.ValueOf(data).IsValid() {
		data = ""
	}
	if msg == "ref" {
		msg = CodeString(code)
	}

	formatter := new(ApiFormat)
	formatter.Code = code
	formatter.Msg = msg
	formatter.Data = data

	return formatter, nil
}
