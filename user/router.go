package user

import (
	"cmdb/common"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

/*
请求参数
*/
type UserRequest struct {
	ID       string `json:"_id"`
	UserName string `json:"user_name" binding:"required"`
	RealName string `json:"real_name" binding:"required"`
	Mobile   string `json:"mobile" binding:"required,check_mobile"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Team     string `json:"team_id"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

/*
登录注册
*/
func UserLogin(router *gin.RouterGroup) {
	router.POST("/login", loginUser)
}

/*
请求路径
*/
func UserEndpoints(router *gin.RouterGroup) {
	CreateUserIndex()
	router.POST("/", createUser)
	router.GET("/", getUser)
	router.GET("/list", getUserList)
	router.PUT("/", updateUser)
	router.DELETE("/", deleteUser)
}

/*
建立用户
*/
func createUser(c *gin.Context) {
	// 请求处理
	var params UserRequest
	var response common.Response

	if err := c.ShouldBindJSON(&params); err != nil {
		response.Code, response.Message = 1001, err.Error()
		c.JSON(200, response)
		return
	}

	// 数据库处理
	document := User{
		ID:       primitive.NewObjectID(),
		UserName: params.UserName,
		RealName: params.RealName,
		Email:    params.Email,
		Mobile:   params.Mobile,
		Password: common.SetPassword(params.Password),
		CreateAt: time.Now().Local().Unix(),
	}

	if len(params.Team) != 0 {
		team_id, _ := primitive.ObjectIDFromHex(params.Team)
		document.Team = team_id
	}

	id, err := UserModel.Mgo.InsertOne(document)
	if err != nil {
		response.Code, response.Message = 2001, err.Error()
		c.JSON(200, response)
		return
	}

	// 响应处理
	response.Code, response.Message = 0, "用户创建成功"
	response.Data = map[string]string{"_id": id.Hex()}
	c.JSON(200, response)
}

/*
获取用户
*/
func getUser(c *gin.Context) {
	// 请求处理
	var response common.Response
	var result *User

	param := c.Query("_id")
	id, _ := primitive.ObjectIDFromHex(param)

	// 数据库处理
	filter := bson.D{{"_id", id}}
	if err := UserModel.Mgo.GetByField(&result, filter); err != nil {
		response.Code, response.Message = 2001, err.Error()
		c.JSON(200, response)
		return
	}

	// 响应处理
	response.Code, response.Message = 0, "用户获取成功"
	response.Data = result
	c.JSON(200, response)
}

/*
更新用户
*/
func updateUser(c *gin.Context) {
	// 请求处理
	var params UserRequest
	var response common.Response
	var result *User

	if err := c.ShouldBindJSON(&params); err != nil {
		response.Code, response.Message = 1001, err.Error()
		c.JSON(200, response)
		return
	}

	// 数据库处理
	id, _ := primitive.ObjectIDFromHex(params.ID)
	team_id, _ := primitive.ObjectIDFromHex(params.Team)

	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"real_name", params.RealName},
			{"email", params.Email},
			{"mobile", params.Mobile},
			{"update_at", time.Now().Local().Unix()},
			{"team_id", team_id},
		},
		},
	}

	if err := UserModel.Mgo.UpdateByField(&result, filter, update); err != nil {
		response.Code, response.Message = 2001, err.Error()
		c.JSON(200, response)
		return
	}

	// 响应处理
	response.Code, response.Message = 0, "用户更新成功"
	response.Data = result
	c.JSON(200, response)
}

/*
删除用户
*/
func deleteUser(c *gin.Context) {
	// 请求处理
	var response common.Response

	param := c.Query("_id")
	id, _ := primitive.ObjectIDFromHex(param)

	// 数据库处理
	filter := bson.D{{"_id", id}}
	count, err := UserModel.Mgo.DeleteByField(filter)
	if err != nil {
		response.Code, response.Message = 2001, err.Error()
		c.JSON(200, response)
		return
	}

	// 响应处理
	if count == 0 {
		response.Code, response.Message = 1002, "用户未找到，无法删除"
		c.JSON(200, response)
		return
	}

	response.Code, response.Message = 0, "用户删除成功"
	response.Data = map[string]string{"_id": id.Hex()}
	c.JSON(200, response)
}

/*
获取用户列表
*/
func getUserList(c *gin.Context) {
	// 请求处理
	var response common.ResponseList

	index := c.DefaultQuery("index", "1")
	limit := c.DefaultQuery("limit", "20")
	sorts := c.QueryArray("sort")
	filter := c.DefaultQuery("filter", "")

	pageIndex, _ := strconv.ParseInt(index, 10, 64)
	pageLimit, _ := strconv.ParseInt(limit, 10, 64)

	var filters bson.D
	if filter == "" {
		filters = bson.D{}
	} else {
		filters = bson.D{
			{"$or", bson.A{
				bson.D{{"email", bson.D{{"$regex", filter}}}},
				bson.D{{"mobile", bson.D{{"$regex", filter}}}},
				bson.D{{"real_name", bson.D{{"$regex", filter}}}},
				bson.D{{"user_name", bson.D{{"$regex", filter}}}},
				bson.D{{"team_name", bson.D{{"$regex", filter}}}},
			}},
		}
	}

	// 数据库处理
	matchStage := bson.D{{"$match", filters}}

	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "team"},
			{"localField", "team_id"},
			{"foreignField", "_id"},
			{"as", "from_team"},
		}},
	}

	replaceRootStage := bson.D{
		{"$replaceRoot", bson.D{
			{"newRoot", bson.D{
				{"$mergeObjects", bson.A{bson.D{{"$arrayElemAt", bson.A{"$from_team", 0}}}, "$$ROOT"}},
			}},
		}},
	}

	projectStage := bson.D{
		{"$project", bson.D{
			{"from_team", 0},
		}},
	}

	pipeline := mongo.Pipeline{lookupStage, replaceRootStage, projectStage, matchStage}

	list, err := UserModel.Mgo.GetList(pageIndex, pageLimit, sorts, filters, pipeline)
	if err != nil {
		response.Code, response.Message = 2001, err.Error()
		c.JSON(200, response)
		return
	}

	// 响应处理
	if list.Total == 0 {
		response.Code, response.Message = 1001, "没有找到数据"
	} else {
		response.Code, response.Message = 0, "用户列表获取成功"
	}
	response.Data = list
	c.JSON(200, response)
}

/*
用户登录
*/
func loginUser(c *gin.Context) {
	// 请求处理
	var params LoginRequest
	var response common.Response
	var result *User

	if err := c.ShouldBindJSON(&params); err != nil {
		response.Code, response.Message = 1001, err.Error()
		c.JSON(200, response)
		return
	}

	// 数据库处理
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"user_name", params.Login}},
			bson.D{{"email", params.Login}},
			bson.D{{"mobile", params.Login}},
		}},
	}

	if err := UserModel.Mgo.GetByField(&result, filter); err != nil {
		response.Code, response.Message = 2001, err.Error()
		c.JSON(200, response)
		return
	}

	// 校验密码
	if common.VerifyPassword(params.Password, result.Password) != nil {
		response.Code, response.Message = 2001, "密码验证失败"
	}

	// 生成token
	_id := result.ID.Hex()
	access_token, refresh_token := common.GenerateToken(_id)

	// 写入redis
	_ = common.RDB.Set(_id+"."+access_token, access_token, time.Minute*15+time.Second*30)
	_ = common.RDB.Set(_id+"."+refresh_token, refresh_token, time.Minute*30+time.Second*30)

	// 响应处理
	response.Code, response.Message = 0, "用户登录成功"
	response.Data = map[string]string{
		"_id":           _id,
		"access_token":  access_token,
		"refresh_token": refresh_token,
	}
	c.JSON(200, response)
}

/*
用户退出
*/
func logoutUser(c *gin.Context) {}

/*
刷新token
*/
func refreshToken(c *gin.Context) {}

/*
密码修改
*/
func changePassword(c *gin.Context) {}
