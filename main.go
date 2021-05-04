package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	router := gin.Default()
	lists = make(map[uuid.UUID]TalkingList)
	readListFromFile()

	setup_routes(router)

	router.Run()
}
