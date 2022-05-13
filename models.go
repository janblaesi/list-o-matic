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

// TalkingListAttendee represents a person that attends an event
type TalkingListAttendee struct {
	// The given name of the attendee
	GivenName string `json:"given_name" binding:"required"`

	// The surname of the attendee
	SurName string `json:"sur_name" binding:"required"`

	// The degree of the attendee
	Degree string `json:"degree" binding:"required"`

	// The e-mail address of the attendee
	Mail string `json:"mail" binding:"-"`
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

	// List of attendees of this event
	Attendees map[uuid.UUID]TalkingListAttendee `json:"attendees" binding:"-"`

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
