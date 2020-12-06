package user

import (
    "cmdb/common"
    "strconv"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/crypto/bcrypt"

    "github.com/gin-gonic/gin"
)

// 用户密码明文加密
func SetPassword(password string) string {
    bytePassword := []byte(password)
    passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
    return string(passwordHash)
}

// 校验用户密码
func VerifyPassword(passwordHash string, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
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
    router.POST("/logout", logoutUser)
    router.POST("/changepassword", changePassword)
}

/*
建立用户
*/
func createUser(c *gin.Context) {
    // 请求处理
    var params UserRequest
    var response common.Response

    if err := c.ShouldBindJSON(&params); err != nil {
        response.Code, response.Message, response.Error = 1001, "请求参数异常", err.Error()
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
        Password: SetPassword(params.Password),
        CreateAt: time.Now().Local().Unix(),
    }

    if len(params.Team) != 0 {
        team_id, _ := primitive.ObjectIDFromHex(params.Team)
        document.Team = team_id
    }

    id, err := UserModel.Mgo.InsertOne(document)
    if err != nil {
        response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
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
        response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
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
        response.Code, response.Message, response.Error = 1001, "请求参数异常", err.Error()
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
        response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
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
        response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
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
            {"as", "team"},
        }},
    }

    replaceRootStage := bson.D{
        {"$replaceRoot", bson.D{
            {"newRoot", bson.D{
                {"$mergeObjects", bson.A{bson.D{{"$arrayElemAt", bson.A{"$team", 0}}}, "$$ROOT"}},
            }},
        }},
    }

    projectStage := bson.D{
        {"$project", bson.D{
            {"team", 0},
        }},
    }

    pipeline := mongo.Pipeline{lookupStage, replaceRootStage, projectStage, matchStage}

    list, err := UserModel.Mgo.GetList(pageIndex, pageLimit, sorts, filters, pipeline)
    if err != nil {
        response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
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
        response.Code, response.Message, response.Error = 1001, "请求参数异常", err.Error()
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

    // 用户验证
    if err := UserModel.Mgo.GetByField(&result, filter); err != nil {
        response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
        c.JSON(200, response)
        return
    }

    if err := VerifyPassword(result.Password, params.Password); err != nil {
        response.Code, response.Message, response.Error = 2001, "密码验证失败", err.Error()
        c.JSON(200, response)
        return
    }

    // 生成token
    _id := result.ID.Hex()
    token, err := GenerateToken(_id)
    if err != nil {
        response.Code, response.Message, response.Error = 2001, "token生成失败", err.Error()
        c.JSON(200, response)
        return
    }

    // redis处理
    _ = common.RDB.Set(_id+":"+token[0], token[0], time.Second*JWT_ACCESS_TOKEN_EXPIRATION)
    _ = common.RDB.Set(_id+":"+token[1], token[1], time.Second*JWT_REFRESH_TOKEN_EXPIRATION)

    // 响应处理
    response.Code, response.Message = 0, "用户登录成功"
    response.Data = map[string]string{
        "_id":           _id,
        "access_token":  token[0],
        "refresh_token": token[1],
    }
    c.JSON(200, response)
}

/*
用户退出
*/
func logoutUser(c *gin.Context) {
    // 请求处理
    var params RefreshTokenRequest
    var response common.Response

    if err := c.ShouldBindJSON(&params); err != nil {
        response.Code, response.Message, response.Error = 1001, "请求参数异常", err.Error()
        c.JSON(200, response)
        return
    }

    // redis处理
    if _, err := common.RDB.Del(params.ID + ":" + params.AccessToken).Result(); err != nil {
        response.Code, response.Message, response.Error = 1001, "access_token未清除，用户注销失败", err.Error()
        c.JSON(200, response)
        return
    }

    if _, err := common.RDB.Del(params.ID + ":" + params.RefreshToken).Result(); err != nil {
        response.Code, response.Message, response.Error = 1001, "refresh_token未清除，用户注销失败", err.Error()
        c.JSON(200, response)
        return
    }

    // 响应处理
    response.Code, response.Message = 0, "用户注销成功"
    c.JSON(200, response)
}

/*
密码修改
*/
func changePassword(c *gin.Context) {
    // 请求处理
    var params ChangePasswordRequest
    var response common.Response
    var result *User

    if err := c.ShouldBindJSON(&params); err != nil {
        response.Code, response.Message, response.Error = 1001, "请求参数异常", err.Error()
        c.JSON(200, response)
        return
    }

    // 数据库处理
    id, _ := primitive.ObjectIDFromHex(params.ID)
    filter := bson.D{{"_id", id}}

    if err := UserModel.Mgo.GetByField(&result, filter); err != nil {
        response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
        c.JSON(200, response)
        return
    }

    if err := VerifyPassword(result.Password, params.Password); err != nil {
        response.Code, response.Message, response.Error = 2001, "密码验证失败", err.Error()
        c.JSON(200, response)
        return
    }

    // 更新密码
    update := bson.D{
        {"$set", bson.D{
            {"password", SetPassword(params.NewPassword)},
            {"update_at", time.Now().Local().Unix()},
        },
        },
    }

    if err := UserModel.Mgo.UpdateByField(&result, filter, update); err != nil {
        response.Code, response.Message, response.Error = 2001, "数据库处理异常", err.Error()
        c.JSON(200, response)
        return
    }

    // 生成token
    _id := result.ID.Hex()
    token, err := GenerateToken(_id)
    if err != nil {
        response.Code, response.Message, response.Error = 2001, "token生成失败", err.Error()
        c.JSON(200, response)
        return
    }

    // redis处理
    if _, err := common.RDB.Del(params.ID + ":" + params.AccessToken).Result(); err != nil {
        response.Code, response.Message, response.Error = 1001, "access_token未清除", err.Error()
        c.JSON(200, response)
        return
    }

    if _, err := common.RDB.Del(params.ID + ":" + params.RefreshToken).Result(); err != nil {
        response.Code, response.Message, response.Error = 1001, "refresh_token未清除", err.Error()
        c.JSON(200, response)
        return
    }

    _ = common.RDB.Set(_id+":"+token[0], token[0], time.Second*JWT_ACCESS_TOKEN_EXPIRATION)
    _ = common.RDB.Set(_id+":"+token[1], token[1], time.Second*JWT_REFRESH_TOKEN_EXPIRATION)

    // 响应处理
    response.Code, response.Message = 0, "用户密码已更新"
    response.Data = map[string]string{
        "_id":           _id,
        "access_token":  token[0],
        "refresh_token": token[1],
    }
    c.JSON(200, response)
}

/*
刷新token
*/
func refreshToken(c *gin.Context) {
    // 请求处理
    var params RefreshTokenRequest
    var response common.Response

    if err := c.ShouldBindJSON(&params); err != nil {
        response.Code, response.Message, response.Error = 1001, "请求参数异常", err.Error()
        c.JSON(200, response)
        return
    }

    // 生成token
    _id := params.ID
    token, err := GenerateToken(_id)
    if err != nil {
        response.Code, response.Message, response.Error = 2001, "token生成失败", err.Error()
        c.JSON(200, response)
        return
    }

    // redis处理
    if _, err := common.RDB.Del(params.ID + ":" + params.AccessToken).Result(); err != nil {
        response.Code, response.Message, response.Error = 1001, "access_token未清除", err.Error()
        c.JSON(200, response)
        return
    }

    if _, err := common.RDB.Del(params.ID + ":" + params.RefreshToken).Result(); err != nil {
        response.Code, response.Message, response.Error = 1001, "refresh_token未清除", err.Error()
        c.JSON(200, response)
        return
    }

    _ = common.RDB.Set(_id+":"+token[0], token[0], time.Second*JWT_ACCESS_TOKEN_EXPIRATION)
    _ = common.RDB.Set(_id+":"+token[1], token[1], time.Second*JWT_REFRESH_TOKEN_EXPIRATION)

    // 响应处理
    response.Code, response.Message = 0, "用户token已刷新"
    response.Data = map[string]string{
        "_id":           _id,
        "access_token":  token[0],
        "refresh_token": token[1],
    }
    c.JSON(200, response)
}
