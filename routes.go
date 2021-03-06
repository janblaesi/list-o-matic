//     __    _      __        ____        __  ___      __  _
//    / /   (_)____/ /_      / __ \      /  |/  /___ _/ /_(_)____
//   / /   / / ___/ __/_____/ / / /_____/ /|_/ / __ `/ __/ / ___/
//  / /___/ (__  ) /_/_____/ /_/ /_____/ /  / / /_/ / /_/ / /__
// /_____/_/____/\__/      \____/     /_/  /_/\__,_/\__/_/\___/
//
// Copyright 2021-2022 Jan Blaesi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
// THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
// CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package main

import (
	"math"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hako/durafmt"
)

func setupRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	// Retrieve all talking lists currently known to the application
	public.GET("/list", func(context *gin.Context) {
		listsFiltered := make(map[uuid.UUID]TalkingList)

		// Filter lists so only public lists are returned
		for key, list := range lists {
			if list.Visibility == 2 {
				listsFiltered[key] = list
			}
		}

		context.JSON(http.StatusOK, listsFiltered)
	})

	// Retrieve all talking lists currently known to the application
	// In contrast to the public endpoint, this will also retrieve
	// private lists
	protected.GET("/list", func(context *gin.Context) {
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

		// This endpoint may only access lists that are unlisted or public
		if listEntry.Visibility == 0 {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, listEntry)
	})

	// Retrieve a specific talking list
	// In contrast to the public endpoint, this will also retrieve
	// private lists
	protected.GET("/list/:uuid", func(context *gin.Context) {
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

	// Update the visibility of a talking list
	protected.POST("/list/:uuid/visibility", func(context *gin.Context) {
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

		var requestData TalkingListVisibilityUpdate
		if err := context.ShouldBindJSON(&requestData); err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		listEntry.Visibility = requestData.NewVisibility
		lists[listUuid] = listEntry
		dumpListToFile()

		context.Status(http.StatusOK)
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

	// Retrieve all attendees in a specific talking list
	protected.GET("/list/:uuid/attendee", func(context *gin.Context) {
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

		context.JSON(http.StatusOK, listEntry.Attendees)
	})

	// Retrieve a single attendee in a specific talking list
	protected.GET("/list/:uuid/attendee/:attendee_uuid", func(context *gin.Context) {
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

		attendeeUuid, err := uuid.Parse(context.Param("attendee_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		attendeeEntry, entryPresent := listEntry.Attendees[attendeeUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.JSON(http.StatusOK, attendeeEntry)
	})

	// Create an attendee in a specific talking list
	protected.POST("/list/:uuid/attendee", func(context *gin.Context) {
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

		var requestData TalkingListAttendee
		if err := context.ShouldBindJSON(&requestData); err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if listEntry.Attendees == nil {
			listEntry.Attendees = make(map[uuid.UUID]TalkingListAttendee)
		}

		attendeeUuid := uuid.New()
		listEntry.Attendees[attendeeUuid] = requestData
		lists[listUuid] = listEntry
		dumpListToFile()

		context.Status(http.StatusCreated)
	})

	// Delete an attendee from a specific talking list
	protected.DELETE("/list/:uuid/attendee/:attendee_uuid", func(context *gin.Context) {
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

		attendeeUuid, err := uuid.Parse(context.Param("attendee_uuid"))
		if err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		_, entryPresent = listEntry.Attendees[attendeeUuid]
		if !entryPresent {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		delete(listEntry.Attendees, attendeeUuid)
		lists[listUuid] = listEntry
		dumpListToFile()

		context.Status(http.StatusOK)
	})

	// Get a Markdown report of an event that may be converted to user-readable PDF format using pandoc
	protected.GET("/list/:uuid/mdreport", func(context *gin.Context) {
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

		// Calculate the time distribution
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

		// Combine all time distribution data into a structure for use in the template
		type GroupTimeDistribution struct {
			GroupName         string
			NumContributions  uint
			TimeShareAbsolute time.Duration
			TimeShareRelative float64
		}
		groupTimeDistributions := make(map[uuid.UUID]GroupTimeDistribution)
		for uuid := range listEntry.Groups {
			var groupTimeDistribution GroupTimeDistribution

			groupTimeDistribution.GroupName = listEntry.Groups[uuid].Name
			groupTimeDistribution.NumContributions = groupNumberContributions[uuid]
			groupTimeDistribution.TimeShareAbsolute = groupTimeShare[uuid]
			groupTimeDistribution.TimeShareRelative = math.Floor(((groupTimeShare[uuid].Seconds()/totalTime.Seconds())*100)*100) / 100

			groupTimeDistributions[uuid] = groupTimeDistribution
		}

		context.Header("Content-Type", "text/markdown")

		// Fill template with data and write to HTTP stream
		reportTemplate := template.Must(template.New("report.got").Funcs(template.FuncMap{
			"prettyDuration": func(duration time.Duration) string {
				return durafmt.Parse(duration).LimitFirstN(1).String()
			},
			"getGroupName": func(uuid uuid.UUID) string {
				groupEntry, entryPresent := listEntry.Groups[uuid]
				if !entryPresent {
					return ""
				}
				return groupEntry.Name
			},
			"timeDistribution": func(uuid uuid.UUID) GroupTimeDistribution {
				return groupTimeDistributions[uuid]
			},
			"timeNow": time.Now,
		}).ParseFiles("report.got"))

		err = reportTemplate.Execute(context.Writer, listEntry)
		if err != nil {
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		context.Status(http.StatusOK)
	})
}
