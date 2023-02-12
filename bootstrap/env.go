package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
	PostgresUser           string `mapstructure:"POSTGRES_USER"`
	PostgresTestUser       string `mapstructure:"POSTGRES_TEST_USER"`
	PostgresPassword       string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresTestPassword   string `mapstructure:"POSTGRES_TEST_PASSWORD"`
	PostgresDevelopmentDB  string `mapstructure:"POSTGRES_DEVELOPMENT_DB"`
	PostgresTestDB         string `mapstructure:"POSTGRES_TEST_DB"`
	PostgresProductionDB   string `mapstructure:"POSTGRES_PRODUCTION_DB"`
	PostgresData           string `mapstructure:"PG_DATA"`
	PostgresHost           string `mapstructure:"PG_HOST"`
}

func NewEnv(envPath string) *Env {
	env := Env{}
	viper.SetConfigFile(envPath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
