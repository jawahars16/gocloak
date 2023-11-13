package db

import (
	"fmt"

	"github.com/jawahars16/gocloak/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	db gorm.DB
}

func New(config config.DBConfig) (Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_general_ci",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return Database{}, err
	}
	return Database{db: *db}, nil
}

func (d *Database) Save(data interface{}) error {
	return d.db.Save(data).Error
}

func (d *Database) First(dest interface{}, conds ...interface{}) error {
	return d.db.First(dest, conds...).Error
}
