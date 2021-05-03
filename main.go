package main

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TalkingListContribution struct {
	PrettyName string    `json:"pretty_name"`
	GroupName  string    `json:"group_name"`
	Duration   time.Time `json:"talking_time"`
}

type TalkingListApplication struct {
	PrettyName string `json:"pretty_name" binding:"required"`
}

type TalkingListGroup struct {
	PrettyName   string                            `json:"pretty_name" binding:"required"`
	Applications map[string]TalkingListApplication `json:"applications"`
}

type TalkingList struct {
	Name                string                      `json:"name" binding:"required"`
	Groups              map[string]TalkingListGroup `json:"groups"`
	CurrentContribution TalkingListContribution     `json:"current_contribution" binding:"-"`
}

var lists map[uuid.UUID]TalkingList

func dumpListToFile() {
	var raw_bytes bytes.Buffer
	enc := gob.NewEncoder(&raw_bytes)

	if err := enc.Encode(lists); err != nil {
		println(err.Error())
		return
	}

	fh, err := os.Create("talking_lists")
	if err != nil {
		println(err.Error())
		return
	}
	defer fh.Close()

	fh.Write(raw_bytes.Bytes())
}

func readListFromFile() {
	var raw_bytes bytes.Buffer
	dec := gob.NewDecoder(&raw_bytes)

	fh, err := os.Open("talking_lists")
	if err != nil {
		println(err.Error())
		return
	}
	defer fh.Close()

	raw_bytes.ReadFrom(fh)

	if err := dec.Decode(&lists); err != nil {
		println(err.Error())
		return
	}
}

func main() {
	router := gin.Default()
	lists = make(map[uuid.UUID]TalkingList)
	readListFromFile()

	router.GET("/list", func(context *gin.Context) {
		context.JSON(http.StatusOK, lists)
	})

	router.GET("/list/:uuid", func(context *gin.Context) {
		list_uuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		list_entry, entry_present := lists[list_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, list_entry)
	})

	router.POST("/list", func(context *gin.Context) {
		list_uuid := uuid.New()

		var request_data TalkingList
		if err := context.ShouldBindJSON(&request_data); err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		lists[list_uuid] = request_data
		dumpListToFile()

		context.Status(http.StatusCreated)
	})

	router.DELETE("/list/:uuid", func(context *gin.Context) {
		list_uuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		_, entry_present := lists[list_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		delete(lists, list_uuid)
		dumpListToFile()

		context.Status(http.StatusOK)
	})

	router.GET("/list/:uuid/group", func(context *gin.Context) {
		list_uuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		list_entry, entry_present := lists[list_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, list_entry.Groups)
	})

	router.POST("/list/:uuid/group/:name", func(context *gin.Context) {
		list_uuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		list_entry, entry_present := lists[list_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		var request_data TalkingListGroup
		if err := context.ShouldBindJSON(&request_data); err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		if list_entry.Groups == nil {
			list_entry.Groups = make(map[string]TalkingListGroup)
		}

		group_name := context.Param("name")
		list_entry.Groups[group_name] = request_data
		lists[list_uuid] = list_entry
		dumpListToFile()

		context.Status(http.StatusCreated)
	})

	router.DELETE("/list/:uuid/group/:name", func(context *gin.Context) {
		list_uuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		list_entry, entry_present := lists[list_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		group_name := context.Param("name")
		delete(list_entry.Groups, group_name)
		lists[list_uuid] = list_entry
		dumpListToFile()

		context.Status(http.StatusOK)
	})

	router.GET("/list/:uuid/group/:name/application", func(context *gin.Context) {
		list_uuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		list_entry, entry_present := lists[list_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		group_name := context.Param("name")
		group_entry, entry_present := list_entry.Groups[group_name]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, group_entry.Applications)
	})

	router.POST("/list/:uuid/group/:name/application/:application_name", func(context *gin.Context) {
		list_uuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		list_entry, entry_present := lists[list_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		group_name := context.Param("name")
		group_entry, entry_present := list_entry.Groups[group_name]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		var request_data TalkingListApplication
		if err := context.ShouldBindJSON(&request_data); err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		if group_entry.Applications == nil {
			group_entry.Applications = make(map[string]TalkingListApplication)
		}

		application_name := context.Param("application_name")
		group_entry.Applications[application_name] = request_data
		list_entry.Groups[group_name] = group_entry
		lists[list_uuid] = list_entry
		dumpListToFile()

		context.Status(http.StatusCreated)
	})

	router.DELETE("/list/:uuid/group/:name/application/:application_name", func(context *gin.Context) {
		list_uuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		list_entry, entry_present := lists[list_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		group_name := context.Param("name")
		group_entry, entry_present := list_entry.Groups[group_name]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		var request_data TalkingListApplication
		if err := context.ShouldBindJSON(&request_data); err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		application_name := context.Param("application_name")
		delete(group_entry.Applications, application_name)
		list_entry.Groups[group_name] = group_entry
		lists[list_uuid] = list_entry
		dumpListToFile()

		context.Status(http.StatusOK)
	})

	router.Run()
}
