package common

import (
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type Jwt struct {
	key        []byte
	expiration int64
}

type AppClaims struct {
	UserId   uint  `json:"uid"`
	Safe     bool  `json:"sf"`
	Email    bool  `json:"em"`
	Phone    bool  `json:"ph"`
	Ga       bool  `json:"ga"`
	VipLevel uint8 `json:"vl"`
	jwt.StandardClaims
}

func JwtSign(key string, minuteString string, userId uint, safe bool, email bool, phone bool, ga bool, vipLevel uint8) (string, int64, error) {
	jwt, err := NewJwt(key, minuteString)
	if err != nil {
		return "", 0, err
	}
	token, err := jwt.Sign(userId, safe, email, phone, ga, vipLevel)
	if err != nil {
		return "", 0, err
	}
	return token, jwt.expiration * 1000, err
}

func JwtParse(key string, minuteString string, tokenString string) (*AppClaims, error) {
	jwt, err := NewJwt(key, minuteString)
	if err != nil {
		return nil, err
	}
	claims, err := jwt.Parse(tokenString)
	if err != nil {
		return nil, err
	}
	return claims, err
}

func NewJwt(key string, minuteString string) (*Jwt, error) {
	minute, err := strconv.ParseInt(minuteString, 10, 64)
	if err != nil {
		return &Jwt{}, err
	}

	return &Jwt{key: []byte(key), expiration: time.Now().Unix() + minute*60}, nil
}

func (j *Jwt) Sign(userId uint, safe bool, email bool, phone bool, ga bool, vipLevel uint8) (string, error) {
	claims := AppClaims{
		userId,
		safe,
		email,
		phone,
		ga,
		vipLevel,
		jwt.StandardClaims{
			ExpiresAt: j.expiration,
			Issuer:    CompanyName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func (j *Jwt) Parse(tokenString string) (*AppClaims, error) {
	var (
		claims *AppClaims
		err    error
	)
	token, err := jwt.ParseWithClaims(tokenString, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.key, nil
	})

	if err != nil {
		return claims, err
	}

	if claims, ok := token.Claims.(*AppClaims); ok && token.Valid {
		return claims, nil
	} else {
		return claims, err
	}
}
