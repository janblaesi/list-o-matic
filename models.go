// List-O-Matic Talking List Management System
// Copyright (C) 2021 Jan Blaesi
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

	// Visibility of this talking list
	// 0 means this list is private, it can only be seen by administrators
	// 1 means this list is unlisted, it can be seen by administrators and those who have a link
	// 2 means this list is public, it can be seen by everybody
	Visibility int `json:"visibility" binding:"-"`

	// Talking groups that are part of this event
	Groups map[uuid.UUID]TalkingListGroup `json:"groups" binding:"-"`

	// The currently active contribution, if there is one.
	// When no one is talking, CurrentContribution.InProgress will be false.
	CurrentContribution TalkingListContribution `json:"current_contribution" binding:"-"`

	// The list of previous contributions
	PastContributions []TalkingListContribution `json:"past_contributions" binding:"-"`
}

// TalkingListVisibilityUpdate represents a request to change the
// visibility of a talking list
type TalkingListVisibilityUpdate struct {
	// The new visibility of the list
	NewVisibility int `json:"new_visibility"`
}
