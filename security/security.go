package security

import (
	"fmt"
	"hakathon/models"
	"log"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	AccessSecret  = []byte("SUPER_SECRET_ACCESS_KEY")
	RefreshSecret = []byte("SUPER_SECRET_REFRESH_KEY")
)

func GenerateToken(AuthID uint, Role string) (models.TokenDataTransfer, error) {
	var err error
	td := models.TokenDataTransfer{
		RefreshTokenID: uuid.New().String(),
		AtExpires:      time.Now().Add(time.Minute * 15).Unix(),
		RtExpires:      time.Now().Add(time.Hour * 7).Unix(),
	}
	AccessTokenData := models.AccessTokenData{
		AuthID: AuthID,
		Role:   Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.AtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	AccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessTokenData)
	td.AccessToken, err = AccessToken.SignedString(AccessSecret)
	if err != nil {
		log.Println("Error token signed!")
		return models.TokenDataTransfer{}, err
	}
	RefreshTokenData := models.RefreshTokenData{
		RtID:   td.RefreshTokenID,
		AuthID: AuthID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.RtExpires, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	RefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshTokenData)
	td.RefreshToken, err = RefreshToken.SignedString(RefreshSecret)
	if err != nil {
		log.Println("Error token signed!")
		return models.TokenDataTransfer{}, err
	}
	return td, nil
}

func ValidateToken(AccessToken string) (models.AccessTokenData, error) {
	var AcToken models.AccessTokenData
	token, err := jwt.ParseWithClaims(AccessToken, &AcToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%v", t.Header["alg"])
		}
		return AccessSecret, nil
	})

	if !token.Valid {
		return models.AccessTokenData{}, fmt.Errorf("not valid")
	}
	if err != nil {
		return models.AccessTokenData{}, err
	}
	return AcToken, err
}

func ValidateRefreshToken(RefreshToken string) (models.RefreshTokenData, error) {
	var RfTokenData models.RefreshTokenData
	RfToken, err := jwt.ParseWithClaims(RefreshToken, &RfTokenData, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%v", token.Header["alg"])
		}
		return RefreshSecret, nil
	})
	if !RfToken.Valid {
		return models.RefreshTokenData{}, fmt.Errorf("Error!")
	}
	if err != nil {
		return models.RefreshTokenData{}, fmt.Errorf("Error: Error to Parse With Claim!")
	}
	return RfTokenData, nil
}

func HashPassword(password string) (string, error) {
	HashP, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return password, err
	}
	return string(HashP), nil
}

func CheckHashPassword(password string, hashpassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(password))
}
