package models

type Response struct {
	Message string `json:"message"`
}

func Write(response string) Response {
	var res Response
	res.Message = response
	return res
}
