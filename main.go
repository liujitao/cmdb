package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "cmdb/common"
    "cmdb/team"
    "cmdb/user"

    "github.com/gin-gonic/gin"
)

func main() {
    client, err := common.InitMgoClient("mongodb://127.0.0.1:27017", "cmdb", 128)
    if err != nil {
        log.Println(err)
    }
    defer client.Disconnect(context.Background())

    route := gin.Default()

    v1 := route.Group("/api/v1")
    team.TeamRegistration(v1.Group("/team"))
    user.UserRegister(v1.Group("/user"))

    if err := route.Run(":8000"); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
