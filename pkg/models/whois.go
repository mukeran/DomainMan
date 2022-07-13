package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

const (
	TableWhois = "whois"
)

type Whois struct {
	ID               uint      `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	FetchedAt        time.Time `json:"fetchedAt"`
	DomainName       string    `json:"domainName"`
	Raw              string    `json:"raw"`
	Registrant       string    `json:"registrant"`
	RegistrantEmail  string    `json:"registrantEmail"`
	Registrar        string    `json:"registrar"`
	RegistrarIANAID  uint      `json:"registrarIanaID"`
	UpdatedDate      time.Time `json:"updatedDate"`
	RegistrationDate time.Time `json:"registrationDate"`
	ExpirationDate   time.Time `json:"expirationDate"`
	Status           uint      `json:"status"`
	NameServer       []string  `gorm:"-" json:"nameServer"`
	NameServer_      []byte    `gorm:"column:name_server" json:"-"`
	DNSSEC           string    `gorm:"column:dnssec" json:"dnssec"`
	DSData           string    `json:"dsData"`
}

func (w *Whois) BeforeSave(db *gorm.DB) (err error) {
	w.NameServer_, err = json.Marshal(w.NameServer)
	return
}

func (w *Whois) AfterFind(db *gorm.DB) (err error) {
	err = json.Unmarshal(w.NameServer_, &w.NameServer)
	return
}
