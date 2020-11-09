package team

import (
    "cmdb/common"
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"

    "github.com/gin-gonic/gin"
)

/*
请求参数
*/
type TeamRequest struct {
    ID       string `json:"_id"`
    TeamName string `json:"team_name" binding:"required"`
}

/*
请求路径
*/
func TeamRegistration(router *gin.RouterGroup) {
    CreateTeamIndex()
    router.POST("/", createTeam)
}

/*
建立团队
*/
func createTeam(c *gin.Context) {
    // 请求处理
    var params *TeamRequest
    var response common.Response

    if err := c.ShouldBindJSON(&params); err != nil {
        response.Code, response.Message = 1001, err.Error()
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
        response.Code, response.Message = 2001, err.Error()
        c.JSON(200, response)
        return
    }

    // 响应处理
    response.Code, response.Message = 0, "团队创建成功"
    response.Data = map[string]string{"_id": id.Hex()}
    c.JSON(200, response)
}
