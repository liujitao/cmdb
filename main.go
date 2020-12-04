package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"time"

	"cmdb/common"
	"cmdb/team"
	"cmdb/user"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// 初始化数据库
	client, err := common.InitMgoClient("mongodb://127.0.0.1:27017", "cmdb", 128)
	if err != nil {
		log.Println(err)
	}
	defer client.Disconnect(context.Background())

	// 初始化redis
	redis := common.InitRedisClient("127.0.0.1:6379", 0, 10)
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
	// initData(client)

	// 启动服务
	if err := route.Run(":8000"); err != nil {
		panic(err)
	}
}

func initData(client *mongo.Client) {
	var teams []interface{}
	var users []interface{}
	var ids []primitive.ObjectID

	var _id primitive.ObjectID

	db := client.Database("cmdb")

	for _, name := range []string{"未分组", "研发", "测试", "运维"} {
		_id = primitive.NewObjectID()
		ids = append(ids, _id)
		teams = append(teams, team.Team{ID: _id, TeamName: name, CreateAt: time.Now().Local().Unix()})

	}

	if _, err := db.Collection("team").InsertMany(context.Background(), teams); err != nil {
		log.Panicln(err)
	}

	users = append(users, user.User{
		ID:       primitive.NewObjectID(),
		UserName: "admin",
		RealName: "管理员",
		Mobile:   "00000000000",
		Email:    "admin@abc.com",
		Password: "admin",
		CreateAt: time.Now().Local().Unix(),
		Team:     ids[0],
	})

	for i := 10; i < 40; i++ {
		id := strconv.Itoa(i)
		users = append(users, user.User{
			ID:       primitive.NewObjectID(),
			UserName: "user" + id,
			RealName: "用户" + id,
			Mobile:   "139000000" + id,
			Email:    "user" + id + "@abc.com",
			Password: user.SetPassword("123456"),
			CreateAt: time.Now().Local().Unix(),
			Team:     ids[rand.Intn(3)],
		})
	}

	if _, err := db.Collection("user").InsertMany(context.Background(), users); err != nil {
		log.Panicln(err)
	}
}
