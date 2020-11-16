package main

import (
	"context"
	"log"

	"cmdb/common"
	"cmdb/team"
	"cmdb/user"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func main() {
	// 初始化数据库
	client, err := common.InitMgoClient("mongodb://127.0.0.1:27017", "cmdb", 128)
	if err != nil {
		log.Println(err)
	}
	defer client.Disconnect(context.Background())

	// 初始化redis
	redis := common.InitRedisClient("127.0.0.1:6379", 1, 10)
	defer redis.Close()

	// 自定义校验器
	binding.Validator = new(common.DefaultValidator)

	// 初始化gin
	//route := gin.Default()
	route := gin.New()

	// 注册请求路径
	v1 := route.Group("/api/v1")
	user.UserLogin(v1.Group("/user"))

	// 调用认证中间件
	v1.Use(user.AuthMiddleware())
	user.UserEndpoints(v1.Group("/user"))
	team.TeamEndpoints(v1.Group("/team"))

	// 初始化数据

	// 启动服务
	if err := route.Run(":8000"); err != nil {
		panic(err)
	}
}
