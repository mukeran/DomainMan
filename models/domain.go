package models

import (
	"encoding/json"
	"time"
)

const (
	TableDomain = "domain"
)

type Domain struct {
	ID              uint `gorm:"primary_key"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Name            string
	NameServer      []string `gorm:"-"`
	NameServer_     []byte   `gorm:"column:name_server" json:"-"`
	IsUpdatingWhois bool
}

func (d *Domain) BeforeSave() (err error) {
	d.NameServer_, err = json.Marshal(d.NameServer)
	return
}

func (d *Domain) AfterFind() (err error) {
	err = json.Unmarshal(d.NameServer_, &d.NameServer)
	return
}
