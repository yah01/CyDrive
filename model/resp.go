package model

import "encoding/json"

type Resp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func PackResp(status int, msg string, data interface{}) Resp {
	if data == nil {
		return Resp{
			Status:  status,
			Message: msg,
		}
	}

	dataBytes, _ := json.Marshal(data)
	return Resp{
		Status:  status,
		Message: msg,
		Data:    string(dataBytes),
	}
}
