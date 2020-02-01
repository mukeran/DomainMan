package models

import "time"

const (
	TableConfig = "config"
)

type Config struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string `gorm:"primary_key"`
	Value     string
}
