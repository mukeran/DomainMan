package templates

import (
	models2 "DomainMan/pkg/models"
	"gorm.io/gorm"
)

type Initial struct {
}

func (Initial) Execute(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(
			&models2.AccessToken{},
			&models2.Config{},
			&models2.Domain{},
			&models2.Registrar{},
			&models2.Suffix{},
			&models2.Whois{}); err != nil {
			return err
		}
		return nil
	})
}
