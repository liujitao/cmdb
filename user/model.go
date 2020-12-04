package user

import (
	"cmdb/common"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserName string             `bson:"user_name" json:"user_name"`
	RealName string             `bson:"real_name" json:"real_name"`
	Mobile   string             `bson:"mobile" json:"mobile"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	CreateAt int64              `bson:"create_at" json:"create_at"`
	UpdateAt int64              `bson:"update_at" json:"update_at"`
	Team     primitive.ObjectID `bson:"team_id" json:"team_id"`
}

type userModel struct {
	Mgo common.Mgo
}

var UserModel = newUserModel()

// 初始化
func newUserModel() *userModel {
	m := new(userModel)
	m.Mgo.SetCollName("user")
	return m
}

// 创建索引
func CreateUserIndex() {
	UserModel.Mgo.CreateIndex(map[string]int{"user_name": 1}, true)
	UserModel.Mgo.CreateIndex(map[string]int{"create_at": 1}, false)
}
