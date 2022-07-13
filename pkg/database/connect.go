package database

import (
	"DomainMan/pkg/errors"
	"DomainMan/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
)

var (
	DB *gorm.DB
)

func ConnectUsingEnv(autoMigrate bool) error {
	dialect := os.Getenv("DOMAINMAN_DATABASE_DIALECT")
	parameter := os.Getenv("DOMAINMAN_DATABASE_PARAMETER")
	return Connect(dialect, parameter, autoMigrate)
}

func Connect(dialect, parameter string, autoMigrate bool) (err error) {
	var d gorm.Dialector
	switch dialect {
	case "mysql":
		d = mysql.Open(parameter)
	case "sqlite":
		d = sqlite.Open(parameter)
	default:
		return errors.ErrUnsupportedDatabaseDialect
	}
	DB, err = gorm.Open(d, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return
	}
	if autoMigrate {
		err = DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(
				&models.AccessToken{},
				&models.Config{},
				&models.Domain{},
				&models.Registrar{},
				&models.Suffix{},
				&models.Whois{}); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return
		}
	}
	return
}
