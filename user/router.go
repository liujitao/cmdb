package user

import (
	"cmdb/common"

	"github.com/gin-gonic/gin"
)

type UserParams struct {
	UserName string `json:"username"`
	RealName string `json:"realname"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Team     string `json:"team_id"`
}

func UserRegister(router *gin.RouterGroup) {
	router.POST("/", CreateUser)
}

func CreateUser(c *gin.Context) {
	var params UserParams
	var response common.Response

	if err := c.ShouldBindJSON(&params); err != nil {
		response.Code, response.Message = common.RequestAccessDeny, "没有权限建立用户"
		c.JSON(400, response)
		return
	}

	response.Code, response.Message = common.RequestSuccess, "用户建立成功"
	response.Data = common.ResponseData{
		Page:     1,
		PageSize: 1,
		Size:     1,
		Total:    1,
		List:     []interface{}{&params},
	}

	c.JSON(201, response)
	return
}
