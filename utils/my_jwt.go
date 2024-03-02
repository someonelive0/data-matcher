package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyJwt struct {
	SigningKey []byte // secret key
}

// Claims on standard claims
type MyClaims struct {
	ID string // authority id
	jwt.StandardClaims
}

// Token errors
var (
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenNotValidYet = errors.New("token not active")
	ErrTokenMalformed   = errors.New("token error")
	ErrTokenInvalid     = errors.New("token invalid")
)

func NewMyJwt(secret []byte) *MyJwt {
	default_secret := []byte("3e7a4194ce5f73845412f4c59e06c20a") // from data-matcher
	if secret != nil {
		default_secret = secret
	}
	return &MyJwt{SigningKey: default_secret}
}

func (p *MyJwt) CreateToken(c MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(p.SigningKey)
}

func (p *MyJwt) ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return p.SigningKey, nil
		})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if token == nil {
		return nil, ErrTokenInvalid
	}

	//parse token to Claims
	//if mc, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	if c, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return c, nil
	}
	return nil, ErrTokenInvalid
}

// RefreshToken refesg token
func (p *MyJwt) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return p.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if c, ok := token.Claims.(*MyClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		c.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return p.CreateToken(*c)
	}
	return "", ErrTokenInvalid
}
