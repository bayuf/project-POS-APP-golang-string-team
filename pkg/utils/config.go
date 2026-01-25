package utils

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Configuration struct {
	AppName     string
	Port        string
	Debug       bool
	Limit       int
	PathLogging string
	GinMode     string
	DB          DatabaseCofig
}

type DatabaseCofig struct {
	Name     string
	Username string
	Password string
	Host     string
	Port     string
	SSL      string
	MaxConn  int32
	MaxIdle  int
	MaxOpen  int
}

func ReadConfiguration() (*Configuration, error) {
	// get config from env file
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return &Configuration{}, err
	}

	// get config from os variable
	viper.AutomaticEnv()

	// get config from flag
	pflag.Int("port-app", 0, "port for app golang")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	return &Configuration{
		AppName:     viper.GetString("APP_NAME"),
		Port:        viper.GetString("PORT"),
		Debug:       viper.GetBool("DEBUG"),
		Limit:       viper.GetInt("LIMIT"),
		PathLogging: viper.GetString("PATH_LOGGING"),
		GinMode:     viper.GetString("GIN_MODE"),

		DB: DatabaseCofig{
			Name:     viper.GetString("DATABASE_NAME"),
			Username: viper.GetString("DATABASE_USERNAME"),
			Password: viper.GetString("DATABASE_PASSWORD"),
			Host:     viper.GetString("DATABASE_HOST"),
			Port:     viper.GetString("DATABASE_PORT"),
			SSL:      viper.GetString("DATABASE_SSL_MODE"),
			MaxConn:  viper.GetInt32("DATABASE_MAX_CONN"),
			MaxIdle:  viper.GetInt("DATABASE_MAX_IDLE_CONN"),
			MaxOpen:  viper.GetInt("DATABASE_MAX_OPEN_CONN"),
		},
	}, nil

}
