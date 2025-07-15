package repository

import (
	"errors"
	"fmt"
	"gityard-api/database"
	"gityard-api/model"
	"gityard-api/security"
	"gorm.io/gorm"
	"time"
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
	if err := db.Model(&user).Where(&model.User{Email: &email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func CreateUserCredential(userId uint, plainPassword string) (*model.UserCredential, error) {
	db := database.DB

	hashedPassword, err := security.HashPassword(plainPassword)
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
	if err := db.Model(&credential).Where(&model.UserCredential{UserID: userId}).First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &credential, nil
}

func CreateOrUpdateUserRefreshToken(userId uint) (*model.UserRefreshToken, error) {
	db := database.DB

	refreshToken, err := security.GenerateRefreshToken(userId)
	if err != nil {
		return nil, err
	}

	data := model.UserRefreshToken{
		RefreshToken: refreshToken.Body,
		ExpiresAt:    time.Now().Add(refreshToken.ExpiresIn),
	}
	userRefreshToken := model.UserRefreshToken{UserID: userId}

	if err := db.Where(&userRefreshToken).Assign(data).FirstOrCreate(&userRefreshToken).Error; err != nil {
		return nil, err
	}

	return &userRefreshToken, nil
}

func GetUserRefreshTokenById(userId uint) (*model.UserRefreshToken, error) {
	db := database.DB

	var refreshToken model.UserRefreshToken
	if err := db.Model(&refreshToken).Where(&model.UserRefreshToken{UserID: userId}).First(&refreshToken).Error; err != nil {
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

	refreshToken, err := security.GenerateRefreshToken(userId)
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
