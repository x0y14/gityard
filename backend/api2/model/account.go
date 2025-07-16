package model

import "time"

// Handlename はシステム内でユニークなハンドルネームを管理します。
type Handlename struct {
	ID         uint      `gorm:"primaryKey"                                                           json:"id"`
	Handlename string    `gorm:"type:varchar(255);not null;uniqueIndex:uq_idx_handlenames_handlename" json:"handlename"`
	CreatedAt  time.Time `gorm:"default:current_timestamp(3)"                                         json:"created_at"`
}

// Account はユーザーに紐づくアカウント（個人・組織）を表します。
type Account struct {
	ID           uint      `gorm:"primaryKey"                                                 json:"id"`
	UserID       uint      `gorm:"not null"                                                   json:"user_id"`
	HandlenameID *uint     `gorm:"uniqueIndex:uq_idx_accounts_handlename_id"                  json:"handlename_id"` // 退会時にNULLになるためポインタ型
	Kind         int       `gorm:"type:smallint;not null;default:1"                           json:"kind"`
	IsDeleted    bool      `gorm:"type:tinyint(1);not null;default:0"                         json:"is_deleted"`
	CreatedAt    time.Time `gorm:"default:current_timestamp(3)"                               json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)" json:"updated_at"`

	// リレーションシップ
	User              User               `gorm:"foreignKey:UserID;constraint:OnDelete:RESTRICT"`
	Handlename        Handlename         `gorm:"foreignKey:HandlenameID;constraint:OnDelete:RESTRICT"`
	AccountPublicKeys []AccountPublicKey `gorm:"foreignKey:AccountID"`
	AccountProfile    AccountProfile     `gorm:"foreignKey:AccountID"`
	Repositories      []Repository       `gorm:"foreignKey:OwnerAccountID"`
}

// AccountPublicKey はアカウントに紐づくSSH公開鍵を表します。
type AccountPublicKey struct {
	ID          uint      `gorm:"primaryKey"                                                                                                            json:"id"`
	AccountID   uint      `gorm:"not null;uniqueIndex:idx_account_id_fingerprint,priority:1"                                                            json:"account_id"`
	FullKeyText string    `gorm:"type:text;not null"                                                                                                    json:"fullkeytext"`
	Algorithm   string    `gorm:"type:varchar(50);not null"                                                                                             json:"algorithm"`
	Keybody     string    `gorm:"type:text;not null"                                                                                                    json:"keybody"`
	Comment     string    `gorm:"type:varchar(255);not null"                                                                                            json:"comment"`
	Fingerprint string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_account_id_fingerprint,priority:2;index:idx_account_publickeys_fingerprint" json:"fingerprint"`
	CreatedAt   time.Time `gorm:"default:current_timestamp(3)"                                                                                          json:"created_at"`
}

// AccountProfile はアカウントの公開プロフィール情報を表します。
type AccountProfile struct {
	AccountID   uint      `gorm:"primaryKey;autoIncrement:false"                                                       json:"account_id"`
	Displayname string    `gorm:"type:varchar(255);not null;default:'unknown';index:idx_account_profiles_displayname"  json:"displayname"`
	Iconpath    string    `gorm:"type:varchar(255);not null;default:'noimage001'"                                      json:"icon_path"`
	IsPrivate   bool      `gorm:"type:tinyint(1);not null;default:0"                                                   json:"is_private"`
	CreatedAt   time.Time `gorm:"default:current_timestamp(3)"                                                         json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)"                           json:"updated_at"`
}

type AccountKind int

const (
	PersonalAccount AccountKind = iota + 1
	OrganizationAccount
)
