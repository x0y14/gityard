package model

import "time"

// Repository はアカウントが所有するリポジトリを表します。
type Repository struct {
	ID             uint      `gorm:"primaryKey"`
	OwnerAccountID *uint     `gorm:"uniqueIndex:uq_idx_repositoris_owner_name,priority:1"` // 所有者削除でNULLになるためポインタ型
	Name           string    `gorm:"type:varchar(255);not null;uniqueIndex:uq_idx_repositoris_owner_name,priority:2"`
	IsPrivate      bool      `gorm:"type:tinyint(1);not null;default:0"`
	CreatedAt      time.Time `gorm:"default:current_timestamp"`
	UpdatedAt      time.Time `gorm:"default:current_timestamp;onUpdate:current_timestamp"`

	// リレーションシップ
	OwnerAccount Account `gorm:"foreignKey:OwnerAccountID;constraint:OnDelete:SET NULL"`
}
