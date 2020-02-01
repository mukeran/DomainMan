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
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	IANAID      uint `gorm:"column:iana_id"`
	Status      uint
	Website     string
	RDAPBaseURL string
	IsFromIANA  bool
	FetchedAt   time.Time
}
