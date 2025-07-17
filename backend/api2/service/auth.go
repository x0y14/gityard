package service

import (
	"gityard-api/database"
	"gityard-api/model"
	"gityard-api/security"
	"gityard-api/service/repository"
	"gorm.io/gorm"
	"log/slog"
)

func SignUp(email, password, handlename string) (*model.User, *model.RefreshToken, error) {
	db := database.DB

	var user *model.User
	var refreshToken *model.RefreshToken
	err := db.Transaction(func(tx *gorm.DB) error {
		// 登録済みでないかチェック
		userInDB, err := repository.GetUserByEmail(tx, email)
		if err != nil {
			return err
		}
		if userInDB != nil { // email登録済み
			return &ErrRegisteredEmail{Email: email}
		}

		handlenameInDB, err := repository.GetHandleNameByName(tx, handlename)
		if err != nil {
			return err
		}
		if handlenameInDB != nil { // handlename登録済み
			return &ErrRegisteredHandleName{HandleName: handlename}
		}

		// 登録処理
		registeredUser, err := repository.CreateUser(tx, email)
		if err != nil {
			return err
		}
		user = registeredUser

		registeredHandleName, err := repository.CreateHandleName(tx, handlename)
		if err != nil {
			return err
		}

		registeredAccount, err := repository.CreateAccount(
			tx,
			registeredUser.ID,
			registeredHandleName.ID,
			model.PersonalAccount,
		)
		if err != nil {
			return err
		}

		_, err = repository.CreateAccountProfile(
			tx,
			registeredAccount.ID,
			registeredHandleName.Handlename,
			false,
		)
		if err != nil {
			return err
		}

		_, err = repository.CreateUserCredential(tx, registeredUser.ID, password)
		if err != nil {
			return err
		}

		registeredRefreshToken, err := repository.CreateOrUpdateUserRefreshToken(tx, registeredUser.ID)
		if err != nil {
			return err
		}
		refreshToken = registeredRefreshToken

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return user, refreshToken, nil
}

func Login(email, password string) (*model.User, *model.RefreshToken, error) {
	db := database.DB

	var user *model.User
	var refreshToken *model.RefreshToken
	err := db.Transaction(func(tx *gorm.DB) error {
		// credentialはuserIdとしか結びついていないので
		userInDB, err := repository.GetUserByEmail(tx, email)
		if err != nil {
			return err
		}
		if userInDB == nil {
			return &ErrUserNotFound{Email: email}
		}
		user = userInDB

		credInDB, err := repository.GetUserCredentialById(tx, userInDB.ID)
		if err != nil {
			return err
		}
		if credInDB == nil {
			return &ErrCredentialNotFound{UserId: userInDB.ID}
		}

		if !security.VerifyPassword(password, credInDB.HashedPassword) {
			return &ErrPasswordMissMatch{UserId: userInDB.ID}
		}

		registeredRefreshToken, err := repository.CreateOrUpdateUserRefreshToken(tx, userInDB.ID)
		if err != nil {
			return err
		}
		refreshToken = registeredRefreshToken

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return user, refreshToken, nil
}

func Logout(userId uint) error {
	db := database.DB

	err := db.Transaction(func(tx *gorm.DB) error {
		refreshTokenInDB, err := repository.GetUserRefreshTokenById(tx, userId)
		if err != nil {
			return err
		}
		if refreshTokenInDB == nil { // リフレッシュトークンを削除しようとしたけど、そもそもなかった。な〜ぜな〜ぜ
			slog.Warn("user refresh token not found")
		}

		return repository.DeleteUserRefreshToken(tx, userId)
	})
	return err
}

func Refresh(refreshToken string) (*uint, *model.RefreshToken, error) {
	db := database.DB

	var userId *uint
	var newRefreshToken *model.RefreshToken
	err := db.Transaction(func(tx *gorm.DB) error {
		userIdInDB, err := repository.GetUserIdByRefreshToken(tx, refreshToken)
		if err != nil {
			return err
		}
		if userIdInDB == nil {
			return &ErrInvalidRefreshTokenProvided{}
		}
		userId = userIdInDB

		generatedRefreshToken, err := repository.CreateOrUpdateUserRefreshToken(tx, *userIdInDB)
		if err != nil {
			return err
		}
		newRefreshToken = generatedRefreshToken

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return userId, newRefreshToken, nil
}
