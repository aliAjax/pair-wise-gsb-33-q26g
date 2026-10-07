package dto

import "time"

// ReminderCreateRequest is the payload for creating a care reminder.
type ReminderCreateRequest struct {
	PlantSpeciesID uint      `json:"plant_species_id"`
	TaskTitle      string    `json:"task_title" binding:"required,max=255"`
	RemindDate     time.Time `json:"remind_date" binding:"required"`
	Frequency      string    `json:"frequency" binding:"omitempty,max=32"`
}

// ReminderStatusRequest carries the new status for a reminder.
type ReminderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
