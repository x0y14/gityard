package repository

import (
	"errors"
	"gityard-api/model"
	"gorm.io/gorm"
)

func CreateHandleName(db *gorm.DB, name string) (*model.Handlename, error) {
	handlename := new(model.Handlename)
	handlename.Handlename = name

	if err := db.Create(&handlename).Error; err != nil {
		return nil, err
	}

	return handlename, nil
}

func GetHandleNameById(db *gorm.DB, handlenameId uint) (*model.Handlename, error) {
	var handlename model.Handlename
	if err := db.Model(&handlename).Where(&model.Handlename{ID: handlenameId}).First(&handlename).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &handlename, nil
}

func GetHandleNameByName(db *gorm.DB, name string) (*model.Handlename, error) {
	var handlename model.Handlename
	if err := db.Model(handlename).Where(&model.Handlename{Handlename: name}).First(&handlename).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &handlename, nil
}

func CreateAccount(db *gorm.DB, userId uint, handlenameId uint, kind model.AccountKind) (*model.Account, error) {
	account := new(model.Account)
	account.UserID = userId
	account.HandlenameID = &handlenameId
	account.Kind = int(kind)

	if err := db.Create(&account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func GetAccountById(db *gorm.DB, accountId uint) (*model.Account, error) {
	var account model.Account
	if err := db.Model(&account).Where(&model.Account{ID: accountId}).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

func CreateAccountProfile(db *gorm.DB, accountId uint, displayName string, private bool) (*model.AccountProfile, error) {
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

func GetAccountProfileById(db *gorm.DB, accountId uint) (*model.AccountProfile, error) {
	var profile model.AccountProfile
	if err := db.Model(&profile).Where(&model.AccountProfile{AccountID: accountId}).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}
