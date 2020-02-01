package templates

import (
	"DomainMan/models"
	"github.com/jinzhu/gorm"
)

type Initial struct {
}

func (Initial) Execute(db *gorm.DB) error {
	tx := db.Begin()
	defer tx.RollbackUnlessCommitted()
	if v := tx.AutoMigrate(&models.AccessToken{}, &models.Config{}, &models.Domain{}, &models.Registrar{}, &models.Suffix{}, &models.Whois{}); v.Error != nil {
		return v.Error
	}
	if v := tx.Commit(); v.Error != nil {
		return v.Error
	}
	return nil
}
