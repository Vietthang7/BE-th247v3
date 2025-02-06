package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	envFile = ".env"
	//certFile = "token_jwt_key.pem"
	configFile = "config.yaml"
)

type AppConfig struct {
	Database   DatabaseConfig `yaml:"database"`
	Logging    LoggingConfig  `yaml:"logging"`
	ConfigFile string
}

var Database *DatabaseConfig
var Logging *LoggingConfig

//var SSOClient *casdoorsdk.Client
//var SSOProviders *casdoorsdk.Provider

func Setup() {
	err := godotenv.Load(envFile)
	if err != nil {
		logrus.Debug(err)
	}

	var Http = &AppConfig{
		ConfigFile: configFile,
		Database: DatabaseConfig{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			Username: os.Getenv("DATABASE_USERNAME"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			DBName:   os.Getenv("DATABASE_NAME"),
		},
	}

	////config sso sdk
	//certContent, err1 := os.ReadFile(certFile)
	//if err1 != nil {
	//	fmt.Println("Err: certfile Not found. Default certContent empty string")
	//	certContent = []byte("")
	//}
	//certKey := string(certContent)
	//_ = os.Setenv("SSO_CERTIFICATE", certKey)
	//SSOClient = casdoorsdk.NewClient(Config("SSO_END_POINT"), Config("SSO_CLIENT_ID"), Config("SSO_CLIENT_SECRET"),
	//	Config("SSO_CERTIFICATE"), Config("SSO_ORGANIZATION_NAME"), Config("SSO_APPLICATION_NAME"))
	//end sso sdk config
	//Http.Database.Setup()
	Http.Database.Setup()
	fmt.Println("*************** APP DATABASE SETUP FINISHED ***************")
	Http.Logging.Setup()
	Database = &Http.Database
	Logging = &Http.Logging
}

func Config(key string) string {
	value := os.Getenv(key)
	return value
}
