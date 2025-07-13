package model

import "time"

// Account はユーザーに紐づくアカウント（個人・組織）を表します。
type Account struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null"`
	HandlenameID *uint     `gorm:"uniqueIndex:uq_idx_accounts_handlename_id"` // 退会時にNULLになるためポインタ型
	Kind         int       `gorm:"type:smallint;not null;default:1"`
	IsDeleted    bool      `gorm:"type:tinyint(1);not null;default:0"`
	CreatedAt    time.Time `gorm:"default:current_timestamp(3)"`
	UpdatedAt    time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)"`

	// リレーションシップ
	User              User               `gorm:"foreignKey:UserID;constraint:OnDelete:RESTRICT"`
	Handlename        Handlename         `gorm:"foreignKey:HandlenameID;constraint:OnDelete:RESTRICT"`
	AccountPublicKeys []AccountPublicKey `gorm:"foreignKey:AccountID"`
	AccountProfile    AccountProfile     `gorm:"foreignKey:AccountID"`
	Repositories      []Repository       `gorm:"foreignKey:OwnerAccountID"`
}

// AccountPublicKey はアカウントに紐づくSSH公開鍵を表します。
type AccountPublicKey struct {
	ID          uint      `gorm:"primaryKey"`
	AccountID   uint      `gorm:"not null;uniqueIndex:idx_account_id_fingerprint,priority:1"`
	Fulltext    string    `gorm:"type:text;not null"`
	Algorithm   string    `gorm:"type:varchar(50);not null"`
	Keybody     string    `gorm:"type:text;not null"`
	Comment     string    `gorm:"type:varchar(255);not null"`
	Fingerprint string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_account_id_fingerprint,priority:2;index:idx_account_publickeys_fingerprint"`
	CreatedAt   time.Time `gorm:"default:current_timestamp(3)"`
}

// AccountProfile はアカウントの公開プロフィール情報を表します。
type AccountProfile struct {
	AccountID   uint      `gorm:"primaryKey"`
	Displayname string    `gorm:"type:varchar(255);not null;default:'unknown';index:idx_account_profiles_displayname"`
	Iconpath    string    `gorm:"type:varchar(255);not null;default:'noimage001'"`
	IsPrivate   bool      `gorm:"type:tinyint(1);not null;default:0"`
	CreatedAt   time.Time `gorm:"default:current_timestamp(3)"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)"`
}
