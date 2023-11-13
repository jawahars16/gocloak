package user

import (
	"log/slog"
	"time"
)

type User struct {
	ID        uint `gorm:"primarykey autoIncrement true not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string
	Password  string
}

type Manager struct {
	db     db
	crypto crypto
}

type crypto interface {
	GenerateFromPassword(password string) (string, error)
}

type db interface {
	Save(data interface{}) error
}

func NewManager(db db, crypto crypto) Manager {
	return Manager{
		db:     db,
		crypto: crypto,
	}
}

func (m *Manager) Add(usr User) error {
	hashedPassword, err := m.crypto.GenerateFromPassword(usr.Password)
	if err != nil {
		slog.Error("Hash password error", err)
		return err
	}

	usr.Password = hashedPassword
	err = m.db.Save(&usr)
	if err != nil {
		slog.Error(err.Error(), "user.email", usr.Email)
	}
	return err
}
