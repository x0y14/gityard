package model

import "time"

type Token struct {
	Type      string        `json:"token_type"`
	ExpiresIn time.Duration `json:"expires_in"`
}

type AccessToken struct {
	Body string `json:"access_token"`
	Token
}

type RefreshToken struct {
	Body string `json:"refresh_token"`
	Token
}
