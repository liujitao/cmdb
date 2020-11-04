package user

import (
	"cmdb/common"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
请求路径
*/
func UserRegister(router *gin.RouterGroup) {
	CreateUserIndex()
	router.POST("/", createUser)
	/*
		router.GET("/", GetUser)
		router.GET("/list", GetUserList)
		router.PUT("/", UpdateUser)
		router.DELETE("/", DeleteUser)
	*/
}

/*
建立用户
*/
func createUser(c *gin.Context) {
	// 请求处理
	var params *UserRequest
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
	var params *UserRequest
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

	// 管理员修改用户密码

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
	//var response common.Response

	/*
		index := c.DefaultQuery("index", "1")
		limit := c.DefaultQuery("limit", "30")

		sort := c.QueryMap("sort")

		team_id, _ := primitive.ObjectIDFromHex("435343fe345fsfsfsf")
		matchStage := bson.D{
			{"$match", bson.D{
				{"$or", bson.A{
					bson.D{{"team_id", team_id}},
				}},
			}},
		}

		pageIndex, _ := strconv.Atoi(index)
		pageLimit, _ := strconv.Atoi(limit)

		limitStage := bson.D{{"$limit", pageLimit}}

		sortStage := bson.D{
			{"$sort", bson.D{
				{},
			}},
		}

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
					{"mergeObjects", bson.A{bson.D{{"$arrayElemAt", bson.A{"$fromteam", 0}}}, "$$ROOT"}},
				}},
			}},
		}

		projectStage := bson.D{
			{"$project", bson.D{
				{"from_team", 0},
				{"password", 0},
				{"team_id", 0},
			}},
		}

		pipeline := mongo.Pipeline{lookupStage, replaceRootStage, projectStage, limitStage, sortStage}
		filter := bson.M{}

		// 数据库处理
		err := UserModel.Mgo.GetList(filter, pipeline)
		if err != nil {
			c.JSON(200, gin.H{"welcome": "bad"})
			return
		}
	*/

	// 响应处理
	c.JSON(200, gin.H{"welcome": "ok"})
}
