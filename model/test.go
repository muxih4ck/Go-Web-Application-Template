package model

type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Token  `json:"data"`
}
