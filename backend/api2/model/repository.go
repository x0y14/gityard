package model

import "time"

// Repository はアカウントが所有するリポジトリを表します。
type Repository struct {
	ID             uint      `gorm:"primaryKey"                                                                      json:"id"`
	OwnerAccountID *uint     `gorm:"uniqueIndex:uq_idx_repositoris_owner_name,priority:1"                            json:"owner_account_id"` // 所有者削除でNULLになるためポインタ型
	Name           string    `gorm:"type:varchar(255);not null;uniqueIndex:uq_idx_repositories_owner_name,priority:2" json:"name"`
	IsPrivate      bool      `gorm:"type:tinyint(1);not null;default:0"                                              json:"is_private"`
	CreatedAt      time.Time `gorm:"default:current_timestamp(3)"                                                    json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:current_timestamp(3);onUpdate:current_timestamp(3)"                      json:"updated_at"`

	// リレーションシップ
	OwnerAccount Account `gorm:"foreignKey:OwnerAccountID;constraint:OnDelete:SET NULL"`
}
