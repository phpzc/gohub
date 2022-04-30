package main

import (
	"fmt"
	"gohub/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {
	//new gin
	router := gin.New()

	bootstrap.SetRoute(router)

	err := router.Run(":3000")

	if err != nil {
		fmt.Println(err.Error())
	}
}
