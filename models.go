package main

import (
	"time"

	"github.com/google/uuid"
)

type TalkingListApplication struct {
	Name string `json:"name" binding:"required"`
}

type TalkingListGroup struct {
	Name         string                               `json:"name" binding:"required"`
	Applications map[uuid.UUID]TalkingListApplication `json:"applications" binding:"-"`
}

type TalkingListContribution struct {
	InProgress  bool                   `json:"in_progress" binding:"-"`
	Application TalkingListApplication `json:"application" binding:"-"`
	GroupUuid   uuid.UUID              `json:"group_uuid" binding:"-"`
	StartTime   time.Time              `json:"start_time" binding:"-"`
	EndTime     time.Time              `json:"end_time" binding:"-"`
	Duration    time.Duration          `json:"duration" binding:"-"`
}

type TalkingList struct {
	Name                string                         `json:"name" binding:"required"`
	Groups              map[uuid.UUID]TalkingListGroup `json:"groups" binding:"-"`
	CurrentContribution TalkingListContribution        `json:"current_contribution" binding:"-"`
	PastContributions   []TalkingListContribution      `json:"past_contributions" binding:"-"`
}

type User struct {
	Username string
	IsAdmin  bool
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required`
	Password string `form:"password" json:"password" binding:"required`
}
