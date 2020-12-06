package team

import (
    "cmdb/common"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

/*
请求参数
*/
type TeamRequest struct {
    ID       string `json:"_id"`
    TeamName string `json:"team_name" binding:"required"`
}

/*
团队
*/

type Team struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
    TeamName string             `bson:"team_name" json:"team_name"`
    CreateAt int64              `bson:"create_at" json:"create_at"`
    UpdateAt int64              `bson:"update_at" json:"update_at"`
}

type teamModel struct {
    Mgo common.Mgo
}

var TeamModel = newTeamModel()

// 初始化
func newTeamModel() *teamModel {
    m := new(teamModel)
    m.Mgo.SetCollName("team")
    return m
}

// 创建索引
func CreateTeamIndex() {
    TeamModel.Mgo.CreateIndex(map[string]int{"team_name": 1}, true)
    TeamModel.Mgo.CreateIndex(map[string]int{"create_at": 1}, false)
}
