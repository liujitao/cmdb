package user

import (
	"cmdb/common"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response common.Response

		// 获取token
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			response.Code, response.Message = 2001, "用户请求非法"
			c.JSON(403, response)
			c.Abort()
			return
		}

		// 解析token
		access_token := strings.Split(token, " ")[1]
		_id := common.ParseToken(access_token)

		// 判断redis中token是否存在
		key := _id + "." + access_token
		_, err := common.RDB.Get(key).Result()
		if err != nil {
			response.Code, response.Message = 2001, "用户token不存在"
			c.JSON(403, response)
			c.Abort()
			return
		}

		// 判断redis中token是否过期
		if ok, _ := common.RDB.Expire(key, time.Minute*15).Result(); ok {
			response.Code, response.Message = 2001, "用户token已过期"
			c.JSON(403, response)
			c.Abort()
			return
		}

		// 生成token
		access_token, refresh_token := common.GenerateToken(_id)

		// 写入redis
		_ = common.RDB.Set(_id+"."+access_token, access_token, time.Minute*15+time.Second*30)
		_ = common.RDB.Set(_id+"."+refresh_token, refresh_token, time.Minute*30+time.Second*30)

		c.Next()
	}
}
