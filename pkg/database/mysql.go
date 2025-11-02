// pkg/database/mysql.go
package database

import (
	"fmt"
	"mini-project-ostore/internal/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func NewMySQLConnection(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Store{},
		&domain.Address{},
		&domain.Category{},
		&domain.Product{},
		&domain.Transaction{},
		&domain.TransactionItem{},
		&domain.ProductLog{},
	)
}
