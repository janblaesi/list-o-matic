package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	router := gin.Default()
	lists = make(map[uuid.UUID]TalkingList)
	readListFromFile()

	setupRoutes(router)

	if err := router.Run(); err != nil {
		println("Failed to start Web Server: ", err.Error())
		return
	}
}
