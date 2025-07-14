package secutiry

import (
	"fmt"
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
	claims["sub"] = strconv.Itoa(int(userId))
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
	claims["sub"] = strconv.Itoa(int(userId))
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

func VerifyAccessToken(accessToken string) (uint, bool) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.Config("SECRET")), nil
	})

	if err != nil || !token.Valid {
		return 0, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	// 用途チェック
	if kind, ok := claims["kind"]; !ok || kind != "access_token" {
		return 0, false
	}

	userIdStr, ok := claims["sub"]
	if !ok {
		return 0, false
	}
	userId64, err := strconv.ParseInt(fmt.Sprintf("%s", userIdStr), 10, 64)
	if err != nil {
		return 0, false
	}

	return uint(userId64), true
}

func VerifyRefreshToken(refreshToken string) (uint, bool) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.Config("SECRET")), nil
	})

	if err != nil || !token.Valid {
		return 0, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}

	// 用途チェック
	if kind, ok := claims["kind"]; !ok || kind != "refresh_token" {
		return 0, false
	}

	userIdStr, ok := claims["sub"]
	if !ok {
		return 0, false
	}
	userId64, err := strconv.ParseInt(fmt.Sprintf("%s", userIdStr), 10, 64)
	if err != nil {
		return 0, false
	}

	return uint(userId64), true
}
