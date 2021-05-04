package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setup_routes(router *gin.Engine) {
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

	router.GET("/list/:uuid/group/:group_uuid", func(context *gin.Context) {
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

		group_uuid, err := uuid.Parse(context.Param("group_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		group_entry, entry_present := list_entry.Groups[group_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, group_entry)
	})

	router.POST("/list/:uuid/group", func(context *gin.Context) {
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
			list_entry.Groups = make(map[uuid.UUID]TalkingListGroup)
		}

		group_uuid := uuid.New()
		list_entry.Groups[group_uuid] = request_data
		lists[list_uuid] = list_entry
		dumpListToFile()

		context.Status(http.StatusCreated)
	})

	router.DELETE("/list/:uuid/group/:group_uuid", func(context *gin.Context) {
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

		group_uuid, err := uuid.Parse(context.Param("group_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		_, entry_present = list_entry.Groups[group_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		delete(list_entry.Groups, group_uuid)
		lists[list_uuid] = list_entry
		dumpListToFile()

		context.Status(http.StatusOK)
	})

	router.GET("/list/:uuid/group/:group_uuid/application", func(context *gin.Context) {
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

		group_uuid, err := uuid.Parse(context.Param("group_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		group_entry, entry_present := list_entry.Groups[group_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, group_entry.Applications)
	})

	router.POST("/list/:uuid/group/:group_uuid/application", func(context *gin.Context) {
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

		group_uuid, err := uuid.Parse(context.Param("group_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		group_entry, entry_present := list_entry.Groups[group_uuid]
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
			group_entry.Applications = make(map[uuid.UUID]TalkingListApplication)
		}

		application_uuid := uuid.New()
		group_entry.Applications[application_uuid] = request_data
		list_entry.Groups[group_uuid] = group_entry
		lists[list_uuid] = list_entry
		dumpListToFile()

		context.Status(http.StatusCreated)
	})

	router.DELETE("/list/:uuid/group/:group_uuid/application/:application_uuid", func(context *gin.Context) {
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

		group_uuid, err := uuid.Parse(context.Param("group_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		group_entry, entry_present := list_entry.Groups[group_uuid]
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

		application_uuid, err := uuid.Parse(context.Param("application_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		delete(group_entry.Applications, application_uuid)
		list_entry.Groups[group_uuid] = group_entry
		lists[list_uuid] = list_entry
		dumpListToFile()

		context.Status(http.StatusOK)
	})

	router.GET("/list/:uuid/start_contribution", func(context *gin.Context) {
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

		group_uuid, err := uuid.Parse(context.Query("group"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		group_entry, entry_present := list_entry.Groups[group_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		application_uuid, err := uuid.Parse(context.Query("application"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		application_entry, entry_present := group_entry.Applications[application_uuid]
		if !entry_present {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		now := time.Now()

		if list_entry.CurrentContribution.InProgress {
			prev_contribution := list_entry.CurrentContribution
			prev_contribution.EndTime = now
			prev_contribution.Duration = prev_contribution.EndTime.Sub(prev_contribution.StartTime)
			prev_contribution.InProgress = false
			list_entry.PastContributions = append(list_entry.PastContributions, prev_contribution)
		}

		list_entry.CurrentContribution.StartTime = now
		list_entry.CurrentContribution.Application = application_entry
		list_entry.CurrentContribution.GroupUuid = group_uuid
		list_entry.CurrentContribution.InProgress = true

		delete(group_entry.Applications, application_uuid)
		list_entry.Groups[group_uuid] = group_entry
		lists[list_uuid] = list_entry
		dumpListToFile()
	})

	router.GET("/list/:uuid/stop_contribution", func(context *gin.Context) {
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

		if list_entry.CurrentContribution.InProgress {
			prev_contribution := list_entry.CurrentContribution
			prev_contribution.EndTime = time.Now()
			prev_contribution.Duration = prev_contribution.EndTime.Sub(prev_contribution.StartTime)
			prev_contribution.InProgress = false
			list_entry.PastContributions = append(list_entry.PastContributions, prev_contribution)
		}

		list_entry.CurrentContribution.InProgress = false
		lists[list_uuid] = list_entry
		dumpListToFile()
	})
}
