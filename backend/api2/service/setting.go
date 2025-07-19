package service

import (
	"gityard-api/database"
	"gityard-api/model"
	"gityard-api/security"
	"gityard-api/service/repository"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"log/slog"
	"strings"
)

func RegisterSSHPublicKey(userId uint, keyName string, pubkeyFullText string) (*model.UserPublicKey, error) {
	db := database.DB

	if strings.Contains(pubkeyFullText, "BEGIN") {
		convertedPk, ok := security.ConvertPKCSToOpenSSH(pubkeyFullText)
		if !ok {
			slog.Info("create ssh pk rejected", "reason", "invalid pk")
			return nil, &ErrInvalidPubkeyProvided{}
		}
		pubkeyFullText = convertedPk
	}

	// 正常なデータか検証
	pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubkeyFullText))
	if err != nil {
		slog.Info("create ssh pk rejected", "detail", err)
		return nil, &ErrInvalidPubkeyProvided{}
	}
	fingerprint := ssh.FingerprintSHA256(pk)

	alg, body, comment, err := security.ParseSSHKey(pubkeyFullText)
	if err != nil {
		slog.Info(pubkeyFullText, err)
		slog.Info("create ssh pk rejected", "reason", "invalid format")
		return nil, &ErrInvalidPubkeyProvided{}
	}

	var pubkey *model.UserPublicKey
	err = db.Transaction(func(tx *gorm.DB) error {
		pubkeyInDB, err := repository.GetPubKeyByFingerprint(tx, fingerprint)
		if err != nil {
			return err
		}
		if pubkeyInDB != nil {
			slog.Warn("pubkey duplicates fingerprint", "fingerprint", fingerprint)
			return &ErrDuplicatesPubkeyFingerprint{}
		}

		pubkey, err = repository.CreatePublicKey(tx, userId, keyName, pubkeyFullText, alg, body, comment, fingerprint)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

func GetSSHPublicKeys(userId uint, offset, limit int) ([]model.UserPublicKey, error) {
	db := database.DB

	var pubkeys []model.UserPublicKey
	err := db.Transaction(func(tx *gorm.DB) error {
		// 存在チェック
		user, err := repository.GetUserById(tx, userId)
		if err != nil {
			return err
		}
		if user == nil {
			return &ErrUserNotFound{UserId: userId}
		}

		pks, err := repository.GetPubkeysByUserId(tx, userId, offset, limit)
		if err != nil {
			return err
		}
		pubkeys = pks

		return nil
	})
	if err != nil {
		return nil, err
	}

	return pubkeys, nil
}

func DeleteSSHPublicKeyByFingerprint(userId uint, fingerprint string) error {
	db := database.DB

	return db.Transaction(func(tx *gorm.DB) error {
		user, err := repository.GetUserById(tx, userId)
		if err != nil {
			return err
		}
		if user == nil {
			return &ErrUserNotFound{UserId: userId}
		}

		return repository.DeletePublicKeyByFingerprint(tx, userId, fingerprint)
	})
}
