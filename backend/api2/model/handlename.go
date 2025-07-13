package model

import "time"

// Handlename はシステム内でユニークなハンドルネームを管理します。
type Handlename struct {
	ID         uint      `gorm:"primaryKey"`
	Handlename string    `gorm:"type:varchar(255);not null;uniqueIndex:uq_idx_handlenames_handlename"`
	CreatedAt  time.Time `gorm:"default:current_timestamp(3)"`
}
