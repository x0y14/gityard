package secutiry

import (
	"github.com/golang-jwt/jwt/v5"
	"gityard-api/config"
	"gityard-api/model"
	"strconv"
	"time"
)

func GenerateAccessToken(userId uint) (*model.AccessToken, error) {
	expiresIn := time.Minute * config.AccessTokenActiveDurationMinutes

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = strconv.Itoa(int(userId))
	claims["exp"] = time.Now().Add(expiresIn).Unix()
	claims["kind"] = "access_token"

	body, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return nil, err
	}

	return &model.AccessToken{
		Body: body,
		Token: model.Token{
			Type:      "Bearer",
			ExpiresIn: expiresIn,
		},
	}, nil
}

func GenerateRefreshToken(userId uint) (*model.RefreshToken, error) {
	expiresIn := time.Minute * config.RefreshTokenActiveDurationMinutes

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = strconv.Itoa(int(userId))
	claims["exp"] = time.Now().Add(expiresIn).Unix()
	claims["kind"] = "refresh_token"

	body, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return nil, err
	}

	return &model.RefreshToken{
		Body: body,
		Token: model.Token{
			Type:      "Bearer",
			ExpiresIn: expiresIn,
		},
	}, nil
}
