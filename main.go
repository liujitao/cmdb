package main

import (
	"fmt"
	"os"

	"cmdb/user"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()

	v1 := route.Group("/api/v1")
	user.UserRegister(v1.Group("/user"))

	if err := route.Run(":8080"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
