package db

import (
	"fmt"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jawahars16/gocloak/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	db     gorm.DB
	config config.DBConfig
}

func New(config config.DBConfig) (Database, error) {
	db, err := gorm.Open(mysql.Open(dsn(config)))
	if err != nil {
		return Database{}, err
	}
	return Database{db: *db, config: config}, nil
}

func (d *Database) Save(data interface{}) error {
	return d.db.Save(data).Error
}

func (d *Database) First(dest interface{}, conds ...interface{}) error {
	return d.db.First(dest, conds...).Error
}

func (d *Database) Migrate() error {
	connStr := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%s)/%s?query",
		d.config.User, d.config.Password, d.config.Host, d.config.Port, d.config.Database)

	mgrt, err := migrate.New("file://infra/db/migrations", connStr)
	if err != nil {
		slog.Error("creating migration. Check if the DB config is correct", "msg", err.Error())
		return err
	}

	err = mgrt.Up()
	if err != nil {
		slog.Error("running migrations", "msg", err.Error())
		return err
	}
	return nil
}

func dsn(config config.DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_general_ci",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database)
}
