package models

import "time"

const (
	TableSuffix = "suffix"
)

type Suffix struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Memo        string
	Description string
	WhoisServer string
}
