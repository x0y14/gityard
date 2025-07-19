package model

import "time"

// Repository はアカウントが所有するリポジトリを表します。
type Repository struct {
	ID             uint      `gorm:"column:id;primaryKey"                                                                                        json:"id"`
	OwnerAccountID *uint     `gorm:"column:owner_account_id;uniqueIndex:uq_idx_repositories_owner_account_id_and_name,priority:1"                json:"owner_account_id"` // 所有者削除でNULLになるためポインタ型
	Name           string    `gorm:"column:name;type:varchar(255);not null;uniqueIndex:uq_idx_repositories_owner_account_id_and_name,priority:2" json:"name"`
	IsPrivate      bool      `gorm:"column:is_private;type:tinyint(1);not null;default:0"                                                        json:"is_private"`
	CreatedAt      time.Time `gorm:"column:created_at;default:current_timestamp(3)"                                                              json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;default:current_timestamp(3);onUpdate:current_timestamp(3)"                                json:"updated_at"`

	// リレーションシップ
	OwnerAccount Account `gorm:"foreignKey:OwnerAccountID;constraint:OnDelete:SET NULL"`
}

func (Repository) TableName() string {
	return "repositories"
}
