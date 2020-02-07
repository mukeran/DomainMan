package models

import (
	"encoding/json"
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

func (d *Domain) BeforeSave() (err error) {
	d.NameServer_, err = json.Marshal(d.NameServer)
	return
}

func (d *Domain) AfterFind() (err error) {
	err = json.Unmarshal(d.NameServer_, &d.NameServer)
	return
}
