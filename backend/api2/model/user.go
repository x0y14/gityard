package model

import (
	"time"
)

// User はユーザーの基本情報を表します。認証の主体となります。
type User struct {
	ID        uint      `gorm:"primaryKey"                                                 json:"id"`
	Email     *string   `gorm:"type:varchar(255);uniqueIndex:uq_idx_users_email"           json:"email"` // 退会時にNULLになるためポインタ型
	IsDeleted bool      `gorm:"type:tinyint(1);not null;default:0"                         json:"is_deleted"`
	CreatedAt time.Time `gorm:"default:current_timestamp(3)"                               json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)" json:"updated_at"`

	// リレーションシップ
	UserCredential   UserCredential   `gorm:"foreignKey:UserID"`
	UserRefreshToken UserRefreshToken `gorm:"foreignKey:UserID"`
	UserPublicKeys   []UserPublicKey  `gorm:"foreignKey:UserID"`
	Accounts         []Account        `gorm:"foreignKey:UserID"`
}

// UserCredential はユーザーのパスワード情報を分離して管理します。
type UserCredential struct {
	UserID         uint      `gorm:"primaryKey;autoIncrement:false"                             json:"user_id"`
	HashedPassword string    `gorm:"type:varchar(255);not null"                                 json:"hashed_password"`
	CreatedAt      time.Time `gorm:"default:current_timestamp(3)"                               json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)" json:"updated_at"`
}

// UserRefreshToken はユーザーのリフレッシュトークンを管理します。
type UserRefreshToken struct {
	UserID             uint      `gorm:"primaryKey;autoIncrement:false"                                           json:"user_id"`
	HashedRefreshToken string    `gorm:"type:varchar(255);not null;uniqueIndex:uq_idx_users_hashed_refresh_token" json:"hashed_refresh_token"`
	ExpiresAt          time.Time `gorm:"not null"                                                                 json:"expires_at"`
	CreatedAt          time.Time `gorm:"default:current_timestamp(3)"                                             json:"created_at"`
	UpdatedAt          time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)"               json:"updated_at"`
}

// UserPublicKey はアカウントに紐づくSSH公開鍵を表します。
type UserPublicKey struct {
	ID          uint      `gorm:"primaryKey"                                                                                                            json:"id"`
	UserID      uint      `gorm:"not null;uniqueIndex:idx_user_id_fingerprint,priority:1"                                                               json:"user_id"`
	Name        string    `gorm:"type:varchar(255);not null"                                                                                            json:"name"`
	FullKeyText string    `gorm:"type:text;not null"                                                                                                    json:"fullkeytext"`
	Algorithm   string    `gorm:"type:varchar(50);not null"                                                                                             json:"algorithm"`
	Keybody     string    `gorm:"type:text;not null"                                                                                                    json:"keybody"`
	Comment     string    `gorm:"type:varchar(255);not null"                                                                                            json:"comment"`
	Fingerprint string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_user_id_fingerprint,priority:2;index:idx_user_publickeys_fingerprint" json:"fingerprint"`
	CreatedAt   time.Time `gorm:"default:current_timestamp(3)"                                                                                          json:"created_at"`
}
