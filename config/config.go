package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const ServerEnvironment = "SERVER_ENVIRONMENT"

func BuildConfigFilePath(configFileName string) string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, configFileName)
}

func LoadServerEnvironmentVars() error {
	// envPath := "../../.env" // works only in debug
	envPath := ".env"
	if _, err := os.Stat(envPath); err == nil {
		err := godotenv.Load(envPath)
		if err != nil {
			return err
		}
	}

	viper.AutomaticEnv()
	return nil
}

func GetMysqlConnectionString() string {
	return viper.GetString("MYSQL_CONNECTION_STRING")
}

func GetMysqlUser() string {
	return viper.GetString("MYSQL_USER")
}

func GetMysqlPassword() string {
	return viper.GetString("MYSQL_PASSWORD")
}

func GetServerPort() string {
	return viper.GetString("SERVER_PORT")
}
