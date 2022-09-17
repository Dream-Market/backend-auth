package config

import "github.com/spf13/viper"

type Config struct {
	Port                 string `mapstructure:"PORT"`
	DBHost               string `mapstructure:"DB_HOST"`
	DBPort               string `mapstructure:"DB_PORT"`
	DBUser               string `mapstructure:"DB_USER"`
	DBPassword           string `mapstructure:"DB_PASSWORD"`
	DBName               string `mapstructure:"DB_NAME"`
	DBConnectionInterval int64  `mapstructure:"DB_CONNECTION_INTERVAL"`
	DBConnectionRetries  int64  `mapstructure:"DB_CONNECTION_RETRIES"`
	JWTSecretKey         string `mapstructure:"JWT_SECRET_KEY"`
	ExpirationHours      int64  `mapstructure:"EXPIRATION_HOURS"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./pkg/config/envs")
	viper.SetConfigName("default")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
