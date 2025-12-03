package bun

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type DB struct {
	Write *gorm.DB
	Read  *gorm.DB
}
type Mongo struct {
	Write *mongo.Database
	Read  *mongo.Database
}
type Redis struct {
	Write *redis.Client
	Read  *redis.Client
}
