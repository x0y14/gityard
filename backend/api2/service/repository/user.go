package repository

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"gityard-api/model"
	"gityard-api/security"
	"gorm.io/gorm"
	"time"
)

func CreateUser(db *gorm.DB, email string) (*model.User, error) {
	user := new(model.User)
	user.Email = &email

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserById(db *gorm.DB, userId uint) (*model.User, error) {
	var user model.User
	if err := db.Model(&user).Where(&model.User{ID: userId}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*model.User, error) {
	var user model.User
	if err := db.Model(&user).Where(&model.User{Email: &email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func CreateUserCredential(db *gorm.DB, userId uint, plainPassword string) (*model.UserCredential, error) {
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

func GetUserCredentialById(db *gorm.DB, userId uint) (*model.UserCredential, error) {
	var credential model.UserCredential
	if err := db.Model(&credential).Where(&model.UserCredential{UserID: userId}).First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &credential, nil
}

func CreateOrUpdateUserRefreshToken(db *gorm.DB, userId uint) (*model.RefreshToken, error) {
	refreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	hashedRefreshToken := fmt.Sprintf("%x", sha256.Sum256([]byte(refreshToken.Body)))
	data := model.UserRefreshToken{
		HashedRefreshToken: hashedRefreshToken,
		ExpiresAt:          time.Now().Add(refreshToken.ExpiresIn),
	}
	userRefreshToken := model.UserRefreshToken{UserID: userId}

	if err := db.Where(&userRefreshToken).Assign(data).FirstOrCreate(&userRefreshToken).Error; err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func GetUserRefreshTokenById(db *gorm.DB, userId uint) (*model.UserRefreshToken, error) {
	var refreshToken model.UserRefreshToken
	if err := db.Model(&refreshToken).Where(&model.UserRefreshToken{UserID: userId}).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &refreshToken, nil
}

func GetUserIdAndExpiresAtByRefreshToken(db *gorm.DB, refreshToken string) (*uint, *time.Time, error) {
	hashedRefreshToken := fmt.Sprintf("%x", sha256.Sum256([]byte(refreshToken)))
	var userRefreshToken model.UserRefreshToken
	if err := db.Model(&userRefreshToken).Where(&model.UserRefreshToken{HashedRefreshToken: hashedRefreshToken}).First(&userRefreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return &userRefreshToken.UserID, &userRefreshToken.ExpiresAt, nil
}

func DeleteUserRefreshToken(db *gorm.DB, userId uint) error {
	return db.Delete(&model.UserRefreshToken{}, userId).Error
}

func CreatePublicKey(db *gorm.DB, userId uint, name, fullKeyText, algorithm, keyBody, comment, fingerprint string) (*model.UserPublicKey, error) {
	pubkey := new(model.UserPublicKey)
	pubkey.UserID = userId
	pubkey.Name = name
	pubkey.FullKeyText = fullKeyText
	pubkey.Algorithm = algorithm
	pubkey.Keybody = keyBody
	pubkey.Comment = comment
	pubkey.Fingerprint = fingerprint

	if err := db.Create(&pubkey).Error; err != nil {
		return nil, err
	}

	return pubkey, nil
}

func GetPubKeyByFingerprint(db *gorm.DB, fingerprint string) (*model.UserPublicKey, error) {
	var pubkey model.UserPublicKey
	if err := db.Model(&pubkey).Where(&model.UserPublicKey{Fingerprint: fingerprint}).First(&pubkey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &pubkey, nil
}

func GetPubkeysByUserId(db *gorm.DB, userId uint, offset, limit int) ([]model.UserPublicKey, error) {
	var pubkeys []model.UserPublicKey
	if err := db.Model(&model.UserPublicKey{}).
		Where(&model.UserPublicKey{UserID: userId}).
		Find(&pubkeys).
		Offset(offset).Limit(limit).Error; err != nil {
		return nil, err
	}

	return pubkeys, nil
}

func DeletePublicKeyByFingerprint(db *gorm.DB, userId uint, fingerprint string) error {
	return db.Where(&model.UserPublicKey{UserID: userId, Fingerprint: fingerprint}).Delete(&model.UserPublicKey{}).Error
}
