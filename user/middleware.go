package user

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 加密密匙
const SecretKey = "W3hJwqbX2MnLJn3Lo+ZOXTPgqwYrszfwH9BkrxTxG0o="

func GenerateToken(_id string) (string, string) {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))

	now := time.Now()
	access_token_exp := now.Add(time.Minute * 15).Unix()
	refresh_token_exp := now.Add(time.Minute * 30).Unix()

	// access_token
	jwt_token.Claims = jwt.MapClaims{
		"_id": _id,
		"exp": access_token_exp,
	}
	access_token, _ := jwt_token.SignedString([]byte(SecretKey))

	// refresh_token: 30min
	jwt_token.Claims = jwt.MapClaims{
		"_id": _id,
		"exp": refresh_token_exp,
	}
	refresh_token, _ := jwt_token.SignedString([]byte(SecretKey))

	return access_token, refresh_token
}

func VerifyToken(tokenString string) (string, int64) {
	var _id string
	var exp int64
	claims := jwt.MapClaims{}

	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_id = claims["_id"].(string)
		exp = int64(claims["exp"].(float64))
	}

	return _id, exp
}
