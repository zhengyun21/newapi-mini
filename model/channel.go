package model

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	Name    string `gorm:"column:name"`
	BaseURL string `gorm:"column:base_url"`
	APIKey  string `gorm:"column:api_key"`
	Models  string `gorm:"column:models"`
	Status  int    `gorm:"column:status"`
}
