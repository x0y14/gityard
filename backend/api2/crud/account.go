package crud

import (
	"errors"
	"gityard-api/database"
	"gityard-api/model"
	"gorm.io/gorm"
)

func CreateHandleName(name string) (*model.Handlename, error) {
	db := database.DB

	handlename := new(model.Handlename)
	handlename.Handlename = name

	if err := db.Create(&handlename).Error; err != nil {
		return nil, err
	}

	return handlename, nil
}

func GetHandleNameById(handlenameid uint) (*model.Handlename, error) {
	db := database.DB

	var handlename model.Handlename
	if err := db.Where(&model.Handlename{ID: handlenameid}).First(&handlename).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &handlename, nil
}

func GetHandleNameByName(name string) (*model.Handlename, error) {
	db := database.DB

	var handlename model.Handlename
	if err := db.Where(&model.Handlename{Handlename: name}).First(&handlename).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &handlename, nil
}

func CreateAccount(userId uint, handlenameId uint, kind model.AccountKind) (*model.Account, error) {
	db := database.DB

	account := new(model.Account)
	account.UserID = userId
	account.HandlenameID = &handlenameId
	account.Kind = int(kind)

	if err := db.Create(&account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func GetAccountById(accountId uint) (*model.Account, error) {
	db := database.DB

	var account model.Account
	if err := db.Where(&model.Account{ID: accountId}).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

func CreateAccountProfile(accountId uint, displayName string, private bool) (*model.AccountProfile, error) {
	db := database.DB

	profile := new(model.AccountProfile)
	profile.AccountID = accountId
	profile.Displayname = displayName
	profile.Iconpath = "noimage001"
	profile.IsPrivate = private

	if err := db.Create(&profile).Error; err != nil {
		return nil, err
	}

	return profile, nil
}

func GetAccountProfileById(accountId uint) (*model.AccountProfile, error) {
	db := database.DB

	var profile model.AccountProfile
	if err := db.Where(&model.AccountProfile{AccountID: accountId}).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}
