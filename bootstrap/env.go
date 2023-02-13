package bootstrap

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                      string
	ServerAddress               string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout              int    `mapstructure:"CONTEXT_TIMEOUT"`
	AccessTokenExpiryHour       int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour      int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret           string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret          string `mapstructure:"REFRESH_TOKEN_SECRET"`
	PostgresDevelopmentUser     string `mapstructure:"POSTGRES_DEVELOPMENT_USER"`
	PostgresDevelopmentPassword string `mapstructure:"POSTGRES_DEVELOPMENT_PASSWORD"`
	PostgresDevelopmentDB       string `mapstructure:"POSTGRES_DEVELOPMENT_DB"`
	PostgresTestUser            string `mapstructure:"POSTGRES_TEST_USER"`
	PostgresTestPassword        string `mapstructure:"POSTGRES_TEST_PASSWORD"`
	PostgresTestDB              string `mapstructure:"POSTGRES_TEST_DB"`
	PostgresProductionUser      string `mapstructure:"POSTGRES_PRODUCTION_USER"`
	PostgresProductionPassword  string `mapstructure:"POSTGRES_PRODUCTION_PASSWORD"`
	PostgresProductionDB        string `mapstructure:"POSTGRES_PRODUCTION_DB"`
	PostgresData                string `mapstructure:"PG_DATA"`
	PostgresHost                string `mapstructure:"PG_HOST"`
	DBPortHostDevelopment       string `mapstructure:"DB_PORT_HOST_DEVELOPMENT"`
	DBPortHostTest              string `mapstructure:"DB_PORT_HOST_TEST"`
}

func NewEnv(envPath string, envMode string) *Env {
	env := Env{}
	viper.SetConfigFile(envPath)
	env.AppEnv = envMode

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	return &env
}
