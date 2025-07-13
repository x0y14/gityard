package secutiry

import (
	"github.com/golang-jwt/jwt/v5"
	"gityard-api/config"
	"strconv"
	"time"
)

func GenerateAccessToken(userId uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = strconv.Itoa(int(userId))
	claims["exp"] = time.Now().Add(time.Minute * config.AccessTokenActiveDurationMinutes).Unix()
	claims["kind"] = "access_token"

	return token.SignedString([]byte(config.Config("SECRET")))
}

func GenerateRefreshToken(userId uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = strconv.Itoa(int(userId))
	claims["exp"] = time.Now().Add(time.Minute * config.RefreshTokenActiveDurationMinutes).Unix()
	claims["kind"] = "refresh_token"

	return token.SignedString([]byte(config.Config("SECRET")))
}
