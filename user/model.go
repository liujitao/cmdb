package user

import (
    "cmdb/common"

    "go.mongodb.org/mongo-driver/bson/primitive"
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

type RefreshTokenRequest struct {
    ID           string `json:"_id" binding:"required"`
    AccessToken  string `json:"access_token" binding:"required"`
    RefreshToken string `json:"refresh_token" binding:"required"`
}

type ChangePasswordRequest struct {
    ID           string `json:"_id" binding:"required"`
    Password     string `json:"password" binding:"required"`
    NewPassword  string `json:"new_password" binding:"required"`
    AccessToken  string `json:"access_token" binding:"required"`
    RefreshToken string `json:"refresh_token" binding:"required"`
}

type VerifyCaptchaRequest struct {
    ID   string `json:"id" binding:"required"`
    Code string `json:"code" binding:"required"`
}

/*
用户
*/
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
