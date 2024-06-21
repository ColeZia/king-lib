package token

import (
	"context"

	ke "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/golang-jwt/jwt"
	"gl.king.im/king-lib/framework"
)

func GenerateToken(secret, name string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name,
		//"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

func GenerateTokenWithAccessKey(secret, key string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"access_key": key,
		//"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

func GetTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromServerContext(ctx)

	if !ok {
		return "", ke.New(400, "METADATA_INFO_ERROR", "元信息获取失败！")
	}

	authToken := md.Get(framework.METADATA_KEY_AUTH_TOKEN)

	return authToken, nil
}
