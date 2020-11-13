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
	route := gin.Default()

	// 注册请求路径
	v1 := route.Group("/api/v1")
	user.UserFirst(v1.Group("/user"))

	team.TeamRegistration(v1.Group("/team"))
	user.UserRegistration(v1.Group("/user"))

	// 初始化数据

	// 测试
	/*
		access, refresh := user.GenerateToken("123456")
		fmt.Println(access, refresh)

		_id, exp := user.VerifyToken(access)
		fmt.Println(_id, exp)

		_id, exp = user.VerifyToken(refresh)
		fmt.Println(_id, exp)
	*/

	// 启动服务
	if err := route.Run(":8000"); err != nil {
		panic(err)
	}
}
