package models

import "time"

const (
	TableConfig = "config"
)

type Config struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Key       string    `gorm:"primary_key" json:"key"`
	Value     string    `json:"value"`
}
