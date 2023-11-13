package auth

import (
	"errors"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/jawahars16/gocloak/config"
	"github.com/jawahars16/gocloak/internal/user"
)

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrMismatchedHashAndPassword = errors.New("mismatched hash and password")
)

type tokenGenerator interface {
	NewWithClaims(method jwt.SigningMethod, claims jwt.Claims, opts ...jwt.TokenOption) *jwt.Token
}

type validator interface {
	CompareHashAndPassword(hashedPassword string, password string) error
}
type db interface {
	First(dest interface{}, conds ...interface{}) error
}

type manager struct {
	db             db
	validator      validator
	tokenGenerator tokenGenerator
	config         config.AuthConfig
}

func NewManager(db db, validator validator, tokenGenerator tokenGenerator, config config.AuthConfig) manager {
	return manager{
		db:             db,
		validator:      validator,
		tokenGenerator: tokenGenerator,
		config:         config,
	}
}

func (m *manager) GenerateToken(email string, password string) (string, error) {
	var user user.User
	err := m.db.First(&user, "email = ?", email)
	if err != nil {
		return "", errors.Join(ErrUserNotFound, err)
	}

	err = m.validator.CompareHashAndPassword(user.Password, password)
	if err != nil {
		return "", errors.Join(ErrMismatchedHashAndPassword, err)
	}

	token := m.tokenGenerator.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": user.Email,
	})

	tokenString, err := token.SignedString([]byte(m.config.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
