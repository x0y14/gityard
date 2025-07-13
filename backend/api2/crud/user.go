package crud

import (
	"errors"
	"fmt"
	"gityard-api/database"
	"gityard-api/model"
	"gityard-api/secutiry"
	"gorm.io/gorm"
)

func CreateUser(email string) (*model.User, error) {
	db := database.DB

	user := new(model.User)
	user.Email = &email

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	db := database.DB

	var user model.User
	if err := db.Where(&model.User{Email: &email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func CreateUserCredential(userId uint, plainPassword string) (*model.UserCredential, error) {
	db := database.DB

	hashedPassword, err := secutiry.HashPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	credential := new(model.UserCredential)
	credential.UserID = userId
	credential.HashedPassword = hashedPassword

	if err := db.Create(credential).Error; err != nil {
		return nil, err
	}

	return credential, nil
}

func GetUserCredentialById(userId uint) (*model.UserCredential, error) {
	db := database.DB

	var credential model.UserCredential
	if err := db.Where(&model.UserCredential{UserID: userId}).First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &credential, nil
}

func CreateUserRefreshToken(userId uint) (*model.UserRefreshToken, error) {
	db := database.DB

	refreshToken, err := secutiry.GenerateRefreshToken(userId)
	if err != nil {
		return nil, err
	}

	userRefreshToken := new(model.UserRefreshToken)
	userRefreshToken.UserID = userId
	userRefreshToken.RefreshToken = refreshToken

	if err := db.Create(&userRefreshToken).Error; err != nil {
		return nil, err
	}

	return userRefreshToken, nil
}

func GetUserRefreshTokenById(userId uint) (*model.UserRefreshToken, error) {
	db := database.DB

	var refreshToken model.UserRefreshToken
	if err := db.Where(&model.UserRefreshToken{UserID: userId}).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &refreshToken, nil
}

func UpdateUserRefreshToken(userId uint) (*model.UserRefreshToken, error) {
	db := database.DB

	userRefreshToken, err := GetUserRefreshTokenById(userId)
	if err != nil {
		return nil, err
	}
	if userRefreshToken == nil {
		return nil, fmt.Errorf("refresh token not found")
	}

	refreshToken, err := secutiry.GenerateRefreshToken(userId)
	if err != nil {
		return nil, err
	}

	if err := db.Model(userRefreshToken).Update("refresh_token", refreshToken).Error; err != nil {
		return nil, err
	}

	return userRefreshToken, nil
}

func DeleteUserRefreshToken(userId uint) error {
	db := database.DB

	return db.Delete(&model.UserRefreshToken{}, userId).Error
}
