package model

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	Key         string `gorm:"column:key;uniqueIndex;type:varchar(191)"`
	Name        string `gorm:"column:name"`
	Status      int    `gorm:"column:status"`
	RemainQuota int    `gorm:"column:remain_quota"`
}
