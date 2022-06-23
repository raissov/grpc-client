package models

type Token struct {
	Token string `json:"access_token"`
}

type VerifyToken struct {
	Response Response `json:"response"`
}

type Response struct {
	Public string `json:"user_public_id"`
}
