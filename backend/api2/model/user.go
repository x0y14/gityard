package model

import (
	"time"
)

// User はユーザーの基本情報を表します。認証の主体となります。
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     *string   `gorm:"type:varchar(255);uniqueIndex:uq_idx_users_email"` // 退会時にNULLになるためポインタ型
	IsDeleted bool      `gorm:"type:tinyint(1);not null;default:0"`
	CreatedAt time.Time `gorm:"default:current_timestamp(3)"`
	UpdatedAt time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)"`

	// リレーションシップ
	UserCredential    UserCredential     `gorm:"foreignKey:UserID"`
	UserRefreshTokens []UserRefreshToken `gorm:"foreignKey:UserID"`
	Accounts          []Account          `gorm:"foreignKey:UserID"`
}

// UserCredential はユーザーのパスワード情報を分離して管理します。
type UserCredential struct {
	UserID         uint      `gorm:"primaryKey"`
	HashedPassword string    `gorm:"type:varchar(255);not null"`
	CreatedAt      time.Time `gorm:"default:current_timestamp(3)"`
	UpdatedAt      time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)"`
}

// UserRefreshToken はユーザーのリフレッシュトークンを管理します。
type UserRefreshToken struct {
	UserID       uint      `gorm:"primaryKey"`
	RefreshToken string    `gorm:"type:varchar(255);not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"default:current_timestamp(3)"`
	UpdatedAt    time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)"`
}
