package models

import "time"

const (
	TableAccessToken  = "access_token"
	AccessTokenLength = 32
)

type AccessToken struct {
	ID        uint         `gorm:"primary_key" json:"id"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Name      string       `json:"name"`
	Token     string       `json:"token,omitempty"`
	IsMaster  bool         `json:"isMaster"`
	CanIssue  bool         `json:"canIssue"`
	Issuer    *AccessToken `gorm:"foreignkey:IssuerID" json:"-"`
	IssuerID  uint         `json:"issuerID"`
}
