package model

// ref: https://gorm.io/ja_JP/docs/conventions.html#%E3%83%86%E3%83%BC%E3%83%96%E3%83%AB%E5%90%8D

type Tabler interface {
	TableName() string
}
