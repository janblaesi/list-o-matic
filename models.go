package main

import (
	"time"

	"github.com/google/uuid"
)

// A TalkingListApplication represents an application to talk at an event.
// It belongs to a TalkingListGroup and will be part of a TalkingListContribution once the person
// is allowed to speak.
type TalkingListApplication struct {
	// Name of the person that wants to speak
	Name string `json:"name" binding:"required"`
}

// A TalkingListGroup represents a group of speakers at an event.
type TalkingListGroup struct {
	// Name of the group
	Name string `json:"name" binding:"required"`

	// A list of TalkingListApplications belonging to said group
	Applications map[uuid.UUID]TalkingListApplication `json:"applications" binding:"-"`
}

// TalkingListContribution represents a contribution to the talking list.
// It is created from a TalkingListApplication.
type TalkingListContribution struct {
	// Indicates, if this contribution is currently active.
	InProgress bool `json:"in_progress" binding:"-"`

	// Reference to the TalkingListApplication this contribution consists of
	Application TalkingListApplication `json:"application" binding:"-"`

	// The UUID of the group the Application belongs to
	GroupUuid uuid.UUID `json:"group_uuid" binding:"-"`

	// The time when the contribution started
	StartTime time.Time `json:"start_time" binding:"-"`

	// The time when the contribution ended
	EndTime time.Time `json:"end_time" binding:"-"`

	// After the contribution is finished, this will contain the delta
	// of StartTime and EndTime
	Duration time.Duration `json:"duration" binding:"-"`
}

// TalkingList represents an event, people may talk at.
type TalkingList struct {
	// Name of the event
	Name string `json:"name" binding:"required"`

	// Talking groups that are part of this event
	Groups map[uuid.UUID]TalkingListGroup `json:"groups" binding:"-"`

	// The currently active contribution, if there is one.
	// When no one is talking, CurrentContribution.InProgress will be false.
	CurrentContribution TalkingListContribution `json:"current_contribution" binding:"-"`

	// The list of previous contributions
	PastContributions []TalkingListContribution `json:"past_contributions" binding:"-"`
}
