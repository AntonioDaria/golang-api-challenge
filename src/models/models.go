package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Action struct {
	ID         int        `json:"id"`
	Type       ActionType `json:"type"`
	UserID     int        `json:"userId"`
	TargetUser int        `json:"targetUser,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type ActionType string

const (
	ActionTypeAddContact   ActionType = "ADD_CONTACT"
	ActionTypeEditContact  ActionType = "EDIT_CONTACT"
	ActionTypeReferUser    ActionType = "REFER_USER"
	ActionTypeViewContacts ActionType = "VIEW_CONTACTS"
)
