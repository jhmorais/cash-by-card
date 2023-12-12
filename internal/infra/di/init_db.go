package di

import (
	"fmt"

	"github.com/jhmorais/cash-by-card/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitGormMysqlDB() (*gorm.DB, error) {
	config.LoadServerEnvironmentVars()

	dsn := fmt.Sprintf("%s:%s@%s", config.GetMysqlUser(), config.GetMysqlPassword(), config.GetMysqlConnectionString())
	// dsn := "root:password@tcp(localhost:3306)/database?charset=utf8&parseTime=True&loc=Local"

	mysqlDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "cashby87_cash_by_card.", // schema name
			SingularTable: true,
		},
	})

	if err != nil {
		return nil, err
	}

	// mysqlDb.AutoMigrate(&entities.Client{})

	// sample.DBSeed(mysqlDb)

	return mysqlDb, err
}
