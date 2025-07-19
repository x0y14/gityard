package model

import (
	"time"
)

// User はユーザーの基本情報を表します。認証の主体となります。
type User struct {
	ID        uint      `gorm:"column:id;primaryKey"                                                 json:"id"`
	Email     *string   `gorm:"column:email;type:varchar(255);uniqueIndex:uq_idx_users_email"           json:"email"` // 退会時にNULLになるためポインタ型
	IsDeleted bool      `gorm:"column:is_deleted;type:tinyint(1);not null;default:0"                         json:"is_deleted"`
	CreatedAt time.Time `gorm:"column:created_at;default:current_timestamp(3)"                               json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:current_timestamp(3);onUpdate:current_timestamp(3)" json:"updated_at"`

	// リレーションシップ
	UserCredential   UserCredential   `gorm:"foreignKey:UserID"`
	UserRefreshToken UserRefreshToken `gorm:"foreignKey:UserID"`
	UserPublicKeys   []UserPublicKey  `gorm:"foreignKey:UserID"`
	Accounts         []Account        `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

// UserCredential はユーザーのパスワード情報を分離して管理します。
type UserCredential struct {
	UserID         uint      `gorm:"column:user_id;primaryKey;autoIncrement:false"                             json:"user_id"`
	HashedPassword string    `gorm:"column:hashed_password;type:varchar(255);not null"                                 json:"hashed_password"`
	CreatedAt      time.Time `gorm:"column:created_at;default:current_timestamp(3)"                               json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;default:current_timestamp(3);onUpdate:current_timestamp(3)" json:"updated_at"`
}

func (UserCredential) TableName() string {
	return "user_credentials"
}

// UserRefreshToken はユーザーのリフレッシュトークンを管理します。
type UserRefreshToken struct {
	UserID             uint      `gorm:"column:user_id;primaryKey;autoIncrement:false"                                           json:"user_id"`
	HashedRefreshToken string    `gorm:"column:hashed_refresh_token;type:varchar(255);not null;uniqueIndex:uq_idx_users_hashed_refresh_token" json:"hashed_refresh_token"`
	ExpiresAt          time.Time `gorm:"column:expires_at;not null"                                                                 json:"expires_at"`
	CreatedAt          time.Time `gorm:"column:created_at;default:current_timestamp(3)"                                             json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;default:current_timestamp(3);onUpdate:current_timestamp(3)"               json:"updated_at"`
}

func (UserRefreshToken) TableName() string {
	return "user_refresh_tokens"
}

// UserPublicKey はアカウントに紐づくSSH公開鍵を表します。
type UserPublicKey struct {
	ID          uint      `gorm:"column:id;primaryKey"                                                                                                            json:"id"`
	UserID      uint      `gorm:"column:user_id;not null;uniqueIndex:idx_user_id_fingerprint,priority:1"                                                               json:"user_id"`
	Name        string    `gorm:"column:name;type:varchar(255);not null"                                                                                            json:"name"`
	FullKeyText string    `gorm:"column:fullkeytext;type:text;not null"                                                                                                    json:"fullkeytext"`
	Algorithm   string    `gorm:"column:algorithm;type:varchar(50);not null"                                                                                             json:"algorithm"`
	Keybody     string    `gorm:"column:keybody;type:text;not null"                                                                                                    json:"keybody"`
	Comment     string    `gorm:"column:comment;type:varchar(255);not null"                                                                                            json:"comment"`
	Fingerprint string    `gorm:"column:fingerprint;type:varchar(255);not null;uniqueIndex:idx_user_id_fingerprint,priority:2;index:idx_user_publickeys_fingerprint" json:"fingerprint"`
	CreatedAt   time.Time `gorm:"column:created_at;default:current_timestamp(3)"                                                                                          json:"created_at"`
}

func (UserPublicKey) TableName() string {
	return "user_publickeys"
}
