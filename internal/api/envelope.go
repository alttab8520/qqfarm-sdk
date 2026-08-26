package api

type Reply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func OK(data any) Reply {
	if data == nil {
		data = map[string]any{}
	}
	return Reply{Code: 0, Msg: "ok", Data: data}
}

func Fail(code int, msg string) Reply {
	return Reply{Code: code, Msg: msg, Data: nil}
}
