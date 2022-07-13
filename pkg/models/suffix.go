package models

import "time"

const (
	TableSuffix = "suffix"

	ModeWhois = iota
	ModeWeb
)

type Suffix struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	Memo        string    `json:"memo"`
	Description string    `json:"description"`
	Mode        uint      `json:"mode"`
	WhoisServer string    `json:"whoisServer"`
}
