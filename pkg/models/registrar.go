package models

import "time"

const (
	Reserved = iota
	Accredited
	Terminated
	TableRegistrar = "registrar"
)

var (
	RegistrarStatusFromString = map[string]uint{
		"reserved":   Reserved,
		"accredited": Accredited,
		"terminated": Terminated,
	}
)

type Registrar struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	IANAID      uint      `gorm:"column:iana_id" json:"ianaID"`
	Status      uint      `json:"status"`
	Website     string    `json:"website"`
	RDAPBaseURL string    `json:"rdapBaseUrl"`
	IsFromIANA  bool      `json:"isFromIana"`
	FetchedAt   time.Time `json:"fetchedAt"`
}
