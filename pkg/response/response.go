package response

import (
	"encoding/json"
	"net/http"

	"github.com/gngtwhh/WBlog/pkg/errcode"
)

type responce struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func result(w http.ResponseWriter, httpStatus int, code int, data any, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	resp := responce{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	json.NewEncoder(w).Encode(resp)
}

// Success returns valid result, pass a optional msg string to overwrite default msg
func Success(w http.ResponseWriter, data any, msgs ...string) {
	msg := errcode.GetMsg(errcode.Success)
	if len(msgs) > 0 && msgs[0] != "" {
		msg = msgs[0]
	}
	result(w, http.StatusOK, errcode.Success, data, msg)
}

// Fail response invalid request, pass a optional msg string to overwrite default msg
func Fail(w http.ResponseWriter, code int, msgs ...string) {
	msg := errcode.GetMsg(code)
	if len(msgs) > 0 && msgs[0] != "" {
		msg = msgs[0]
	}
	result(w, http.StatusOK, code, nil, msg)
}
