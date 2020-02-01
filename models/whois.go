package models

import (
	"encoding/json"
	"time"
)

const (
	TableWhois = "whois"
)

type Whois struct {
	ID               uint `gorm:"primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DomainName       string
	Raw              string
	Registrant       string
	RegistrantEmail  string
	Registrar        string
	RegistrarIANAID  uint
	UpdatedDate      time.Time
	RegistrationDate time.Time
	ExpirationDate   time.Time
	Status           uint
	NameServer       []string `gorm:"-"`
	NameServer_      []byte   `gorm:"column:name_server" json:"-"`
	DNSSEC           string   `gorm:"column:dnssec"`
	DSData           string
}

func (w *Whois) BeforeSave() (err error) {
	w.NameServer_, err = json.Marshal(w.NameServer)
	return
}

func (w *Whois) AfterFind() (err error) {
	err = json.Unmarshal(w.NameServer_, &w.NameServer)
	return
}
