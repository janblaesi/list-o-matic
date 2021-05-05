package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	initDb()
	setupRoutes(router)

	if err := router.Run(); err != nil {
		println("Failed to start Web Server: ", err.Error())
		return
	}
}
