package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

var (
	DB *gorm.DB
)

func Connect() (err error) {
	dialect := os.Getenv("DOMAINMAN_DATABASE_DIALECT")
	parameter := os.Getenv("DOMAINMAN_DATABASE_PARAMETER")
	DB, err = gorm.Open(dialect, parameter)
	if err != nil {
		return
	}
	DB.SingularTable(true)
	return
}
