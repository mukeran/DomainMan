package models

import "time"

const (
	TableAccessToken = "access_token"
)

type AccessToken struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Token     string `json:",omitempty"`
	IsMaster  bool
	CanIssue  bool
	Issuer    *AccessToken `gorm:"foreignkey:IssuerID" json:"-"`
	IssuerID  uint
}
