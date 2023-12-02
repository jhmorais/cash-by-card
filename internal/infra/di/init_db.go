package di

import (
	"fmt"

	"github.com/jhmorais/cash-by-card/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGormMysqlDB() (*gorm.DB, error) {
	config.LoadServerEnvironmentVars()

	dsn := fmt.Sprintf("%s:%s@%s", config.GetMysqlUser(), config.GetMysqlPassword(), config.GetMysqlConnectionString())

	mysqlDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// mysqlDb.AutoMigrate(&entities.Device{})

	// sample.DBSeed(mysqlDb)

	return mysqlDb, err
}
