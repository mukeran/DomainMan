package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

const (
	TableDomain = "domain"
)

type Domain struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Name            string    `json:"name"`
	NameServer      []string  `gorm:"-" json:"nameServer"`
	NameServer_     []byte    `gorm:"column:name_server" json:"-"`
	IsUpdatingWhois bool      `json:"isUpdatingWhois"`
}

func (d *Domain) BeforeSave(db *gorm.DB) (err error) {
	d.NameServer_, err = json.Marshal(d.NameServer)
	return
}

func (d *Domain) AfterFind(db *gorm.DB) (err error) {
	err = json.Unmarshal(d.NameServer_, &d.NameServer)
	return
}
