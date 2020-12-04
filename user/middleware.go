package user

import (
	"cmdb/common"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWT
const JWT_SECRET = "W3hJwqbX2MnLJn3Lo+ZOXTPgqwYrszfwH9BkrxTxG0o="
const JWT_ACCESS_TOKEN_EXPIRATION = 60 * 10
const JWT_REFRESH_TOKEN_EXPIRATION = 60 * 60 * 24 * 1

// 生成token pair
func GenerateToken(_id string) ([]string, error) {
	now := time.Now()
	access_exp := now.Add(time.Second * JWT_ACCESS_TOKEN_EXPIRATION).Unix()
	refresh_exp := now.Add(time.Second * JWT_REFRESH_TOKEN_EXPIRATION).Unix()

	// access_token
	access_claims := jwt.MapClaims{
		"_id": _id,
		"exp": access_exp,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, access_claims)
	access_token, err := at.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return nil, err
	}

	// refresh_token
	refresh_claims := jwt.MapClaims{
		"_id": _id,
		"exp": refresh_exp,
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_claims)
	refresh_token, _ := rt.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return nil, err
	}

	return []string{access_token, refresh_token}, err
}

// 解析token
func ParseToken(tokenString string) (string, error) {
	var _id string
	// var exp int64
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_id = claims["_id"].(string)
		// exp = int64(claims["exp"].(float64))
		return _id, err
	}

	return "", err
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response common.Response

		// 获取access_token
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			response.Code, response.Message = 2001, "token异常"
			c.JSON(401, response)
			c.Abort()
			return
		}

		// 解析access_token
		access_token := strings.Split(token, " ")[1]
		_id, err := ParseToken(access_token)
		if err != nil {
			response.Code, response.Message = 2001, "token异常"
			c.JSON(401, response)
			c.Abort()
			return
		}

		// access_token是否过期
		_, err = common.RDB.Get(_id + ":" + access_token).Result()
		//value, err := common.RDB.Get(_id + "" + access_token).Result()
		if err != nil {
			response.Code, response.Message = 2001, "token已过期"
			c.JSON(401, response)
			c.Abort()
			return
		}

		c.Next()
	}
}
