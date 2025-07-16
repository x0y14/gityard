package model

import (
	"time"
)

// User はユーザーの基本情報を表します。認証の主体となります。
type User struct {
	ID        uint      `gorm:"primaryKey"                                                 json:"id"`
	Email     *string   `gorm:"type:varchar(255);uniqueIndex:uq_idx_users_email"           json:"email"` // 退会時にNULLになるためポインタ型
	IsDeleted bool      `gorm:"type:tinyint(1);not null;default:0"                         json:"is_deleted"`
	CreatedAt time.Time `gorm:"default:current_timestamp(3)"                               json:"craeted_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)" json:"updated_at"`

	// リレーションシップ
	UserCredential   UserCredential   `gorm:"foreignKey:UserID"`
	UserRefreshToken UserRefreshToken `gorm:"foreignKey:UserID"`
	Accounts         []Account        `gorm:"foreignKey:UserID"`
}

// UserCredential はユーザーのパスワード情報を分離して管理します。
type UserCredential struct {
	UserID         uint      `gorm:"primaryKey"                                                 json:"user_id"`
	HashedPassword string    `gorm:"type:varchar(255);not null"                                 json:"hashed_password"`
	CreatedAt      time.Time `gorm:"default:current_timestamp(3)"                               json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)" json:"updated_at"`
}

// UserRefreshToken はユーザーのリフレッシュトークンを管理します。
type UserRefreshToken struct {
	UserID             uint      `gorm:"primaryKey"                                                               json:"user_id"`
	HashedRefreshToken string    `gorm:"type:varchar(255);not null;uniqueIndex:uq_idx_users_hashed_refresh_token" json:"hashed_refresh_token"`
	ExpiresAt          time.Time `gorm:"not null"                                                                 json:"expires_at"`
	CreatedAt          time.Time `gorm:"default:current_timestamp(3)"                                             json:"created_at"`
	UpdatedAt          time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)"               json:"updated_at"`
}
