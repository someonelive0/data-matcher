package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// RS256 只需要一个SignKey
type MyJwt struct {
	signKey []byte // secret key
}

// 自定义声明
type MyClaims struct {
	*jwt.StandardClaims
	ID string // authority id
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
	return &MyJwt{signKey: default_secret}
}

func (p *MyJwt) CreateToken(id string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS256"))

	t.Claims = &MyClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 87600).Unix(), // 10 years
			IssuedAt:  time.Now().Unix(),
			Issuer:    "data-matcher", // 签发人
			Subject:   "manage-api",
			Audience:  "idm",
			Id:        id,
		},
		id,
	}

	return t.SignedString(p.signKey)
}

func (p *MyJwt) ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return p.signKey, nil // 回调使用signKey
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
	if c, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return c, nil
	}
	return nil, ErrTokenInvalid
}
