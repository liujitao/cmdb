package team

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
请求路径
*/
func TeamEndpoints(router *gin.RouterGroup) {
	CreateTeamIndex()
	router.POST("/", createTeam)
	router.PUT("/", updateTeam)
	router.DELETE("/", deleteTeam)
	router.GET("/list", getTeamList)
}

/*
建立团队
*/
func createTeam(c *gin.Context) {
	// 请求处理
	var params TeamRequest
	var response common.Response

	if err := c.ShouldBindJSON(&params); err != nil {
		response.Code, response.Message, response.Error = 1001, "请求参数异常", err.Error()
		c.JSON(200, response)
		return
	}

	// 数据库处理
	document := Team{
		ID:       primitive.NewObjectID(),
		TeamName: params.TeamName,
		CreateAt: time.Now().Local().Unix(),
	}

	id, err := TeamModel.Mgo.InsertOne(document)
	if err != nil {
		response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
		c.JSON(200, response)
		return
	}

	// 响应处理
	response.Code, response.Message = 0, "团队创建成功"
	response.Data = map[string]string{"_id": id.Hex()}
	c.JSON(200, response)
}

/*
更新团队
*/
func updateTeam(c *gin.Context) {
	// 请求处理
	var params TeamRequest
	var response common.Response
	var result *Team

	if err := c.ShouldBindJSON(&params); err != nil {
		response.Code, response.Message, response.Error = 1001, "请求参数异常", err.Error()
		c.JSON(200, response)
		return
	}

	// 数据库处理
	id, _ := primitive.ObjectIDFromHex(params.ID)

	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{
			{"team_name", params.TeamName},
			{"update_at", time.Now().Local().Unix()},
		},
		},
	}

	if err := TeamModel.Mgo.UpdateByField(&result, filter, update); err != nil {
		response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
		c.JSON(200, response)
		return
	}

	// 响应处理
	response.Code, response.Message = 0, "团队更新成功"
	response.Data = result
	c.JSON(200, response)
}

/*
删除团队
*/
func deleteTeam(c *gin.Context) {
	// 请求处理
	var response common.Response

	param := c.Query("_id")
	id, _ := primitive.ObjectIDFromHex(param)

	// 数据库处理
	filter := bson.D{{"_id", id}}
	count, err := TeamModel.Mgo.DeleteByField(filter)
	if err != nil {
		response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
		c.JSON(200, response)
		return
	}

	// 响应处理
	if count == 0 {
		response.Code, response.Message = 1002, "团队未找到，无法删除"
		c.JSON(200, response)
		return
	}

	response.Code, response.Message = 0, "团队删除成功"
	response.Data = map[string]string{"_id": id.Hex()}
	c.JSON(200, response)
}

/*
获取团队列表
*/
func getTeamList(c *gin.Context) {
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
				bson.D{{"team_name", bson.D{{"$regex", filter}}}},
			}},
		}
	}

	// 数据库处理
	matchStage := bson.D{{"$match", filters}}

	lookupStage := bson.D{
		{"$lookup", bson.D{
			{"from", "user"},
			{"localField", "_id"},
			{"foreignField", "team_id"},
			{"as", "user"},
		}},
	}

	replaceRootStage := bson.D{
		{"$replaceRoot", bson.D{
			{"newRoot", bson.D{
				{"$mergeObjects", bson.A{bson.D{{"$arrayElemAt", bson.A{"$user", 0}}}, "$$ROOT"}},
			}},
		}},
	}

	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 1},
			{"team_name", 1},
			{"create_at", 1},
			{"update_at", 1},
			{"user_count", bson.D{{"$size", "$user"}}},
			{"user._id", 1},
			{"user.real_name", 1},
		}},
	}

	pipeline := mongo.Pipeline{lookupStage, replaceRootStage, projectStage, matchStage}

	list, err := TeamModel.Mgo.GetList(pageIndex, pageLimit, sorts, filters, pipeline)
	if err != nil {
		response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
		c.JSON(200, response)
		return
	}

	// 响应处理
	if list.Total == 0 {
		response.Code, response.Message = 1001, "没有找到数据"
	} else {
		response.Code, response.Message = 0, "团队列表获取成功"
	}
	response.Data = list
	c.JSON(200, response)
}
