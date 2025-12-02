package bun

import "gorm.io/gorm"

type DB struct {
	Write *gorm.DB
	Read  *gorm.DB
}
