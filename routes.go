package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setupRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	// Retrieve all talking lists currently known to the application
	public.GET("/list", func(context *gin.Context) {
		context.JSON(http.StatusOK, lists)
	})

	// Retrieve a specific talking list
	public.GET("/list/:uuid", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, listEntry)
	})

	// Create a new talking list
	protected.POST("/list", func(context *gin.Context) {
		listUuid := uuid.New()

		var requestData TalkingList
		if err := context.ShouldBindJSON(&requestData); err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		groupUuid := uuid.New()
		var groupData TalkingListGroup
		groupData.Name = "Redner"
		requestData.Groups = make(map[uuid.UUID]TalkingListGroup)
		requestData.Groups[groupUuid] = groupData

		lists[listUuid] = requestData

		dumpListToFile()

		context.Status(http.StatusCreated)
	})

	// Delete a talking list
	protected.DELETE("/list/:uuid", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		_, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		delete(lists, listUuid)
		dumpListToFile()

		context.Status(http.StatusOK)
	})

	// Retrieve all groups in a specific talking list
	public.GET("/list/:uuid/group", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, listEntry.Groups)
	})

	// Retrieve a single group in a specific talking list
	public.GET("/list/:uuid/group/:group_uuid", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		groupUuid, err := uuid.Parse(context.Param("group_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		groupEntry, entryPresent := listEntry.Groups[groupUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, groupEntry)
	})

	// Create a group in a specific talking list
	protected.POST("/list/:uuid/group", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		var requestData TalkingListGroup
		if err := context.ShouldBindJSON(&requestData); err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		if listEntry.Groups == nil {
			listEntry.Groups = make(map[uuid.UUID]TalkingListGroup)
		}

		groupUuid := uuid.New()
		listEntry.Groups[groupUuid] = requestData
		lists[listUuid] = listEntry
		dumpListToFile()

		context.Status(http.StatusCreated)
	})

	// Delete a group from a specific talking list
	protected.DELETE("/list/:uuid/group/:group_uuid", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		groupUuid, err := uuid.Parse(context.Param("group_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		_, entryPresent = listEntry.Groups[groupUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		delete(listEntry.Groups, groupUuid)
		lists[listUuid] = listEntry
		dumpListToFile()

		context.Status(http.StatusOK)
	})

	// Get the time distribution between groups in a specific talking list
	public.GET("/list/:uuid/time_distribution", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		groupTimeShare := make(map[uuid.UUID]time.Duration)
		groupNumberContributions := make(map[uuid.UUID]uint)
		totalTime := time.Duration(0)
		for uuid := range listEntry.Groups {
			groupTimeShare[uuid] = 0
		}
		for _, contribution := range listEntry.PastContributions {
			groupTimeShare[contribution.GroupUuid] += contribution.Duration
			groupNumberContributions[contribution.GroupUuid]++
			totalTime += contribution.Duration
		}

		context.JSON(http.StatusOK, gin.H{
			"time_share":           groupTimeShare,
			"total_time":           totalTime,
			"number_contributions": groupNumberContributions,
		})
	})

	// Reset the list of previous contributions in a specific talking list
	protected.GET("/list/:uuid/reset_past_contributions", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		listEntry.PastContributions = make([]TalkingListContribution, 0)
		lists[listUuid] = listEntry
		dumpListToFile()

		context.Status(http.StatusOK)
	})

	// Get the list of applications in a specific talking group
	public.GET("/list/:uuid/group/:group_uuid/application", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		groupUuid, err := uuid.Parse(context.Param("group_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		groupEntry, entryPresent := listEntry.Groups[groupUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, groupEntry.Applications)
	})

	// Add an application in a specific talking group
	public.POST("/list/:uuid/group/:group_uuid/application", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		groupUuid, err := uuid.Parse(context.Param("group_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		groupEntry, entryPresent := listEntry.Groups[groupUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		var requestData TalkingListApplication
		if err := context.ShouldBindJSON(&requestData); err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		if groupEntry.Applications == nil {
			groupEntry.Applications = make(map[uuid.UUID]TalkingListApplication)
		}

		applicationUuid := uuid.New()
		groupEntry.Applications[applicationUuid] = requestData
		listEntry.Groups[groupUuid] = groupEntry
		lists[listUuid] = listEntry
		dumpListToFile()

		context.JSON(http.StatusCreated, gin.H{
			"uuid": applicationUuid,
		})
	})

	// Delete an application from a talking group
	public.DELETE("/list/:uuid/group/:group_uuid/application/:application_uuid", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		groupUuid, err := uuid.Parse(context.Param("group_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		groupEntry, entryPresent := listEntry.Groups[groupUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		applicationUuid, err := uuid.Parse(context.Param("application_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		delete(groupEntry.Applications, applicationUuid)
		listEntry.Groups[groupUuid] = groupEntry
		lists[listUuid] = listEntry
		dumpListToFile()

		context.Status(http.StatusOK)
	})

	// Start the contribution (from an application)
	protected.GET("/list/:uuid/start_contribution", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		groupUuid, err := uuid.Parse(context.Query("group"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		groupEntry, entryPresent := listEntry.Groups[groupUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		applicationUuid, err := uuid.Parse(context.Query("application"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		applicationEntry, entryPresent := groupEntry.Applications[applicationUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		now := time.Now()

		if listEntry.CurrentContribution.InProgress {
			prevContribution := listEntry.CurrentContribution
			prevContribution.EndTime = now
			prevContribution.Duration = prevContribution.EndTime.Sub(prevContribution.StartTime)
			prevContribution.InProgress = false
			listEntry.PastContributions = append(listEntry.PastContributions, prevContribution)
		}

		listEntry.CurrentContribution.StartTime = now
		listEntry.CurrentContribution.Application = applicationEntry
		listEntry.CurrentContribution.GroupUuid = groupUuid
		listEntry.CurrentContribution.InProgress = true

		delete(groupEntry.Applications, applicationUuid)
		listEntry.Groups[groupUuid] = groupEntry
		lists[listUuid] = listEntry
		dumpListToFile()
	})

	// Stop the current application
	protected.GET("/list/:uuid/stop_contribution", func(context *gin.Context) {
		listUuid, err := uuid.Parse(context.Param("uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			println(err.Error())
			return
		}

		listEntry, entryPresent := lists[listUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		if listEntry.CurrentContribution.InProgress {
			prevContribution := listEntry.CurrentContribution
			prevContribution.EndTime = time.Now()
			prevContribution.Duration = prevContribution.EndTime.Sub(prevContribution.StartTime)
			prevContribution.InProgress = false
			listEntry.PastContributions = append(listEntry.PastContributions, prevContribution)
		}

		listEntry.CurrentContribution.InProgress = false
		lists[listUuid] = listEntry
		dumpListToFile()
	})
}
