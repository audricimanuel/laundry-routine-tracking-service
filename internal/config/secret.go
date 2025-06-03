package config

import "github.com/spf13/viper"

func ViperBind() {
	viper.BindEnv("ENV")

	// Binding Swagger Auth
	viper.BindEnv("SWAGGER_USERNAME")
	viper.BindEnv("SWAGGER_PASSWORD")

	// Binding JWT
	viper.BindEnv("JWT_SECRET")

	// Binding host
	viper.BindEnv("HOST_ADDRESS")
	viper.BindEnv("HOST_PORT")
	viper.BindEnv("HOST_WRITE_TIMEOUT")
	viper.BindEnv("HOST_READ_TIMEOUT")
	viper.BindEnv("HOST_IDLE_TIMEOUT")
	viper.BindEnv("FE_BASE_URL")

	// Binding Database
	viper.BindEnv("POSTGRES_DB_HOST")
	viper.BindEnv("POSTGRES_DB_USER")
	viper.BindEnv("POSTGRES_DB_PASSWORD")
	viper.BindEnv("POSTGRES_DB_NAME")
	viper.BindEnv("POSTGRES_DB_PORT")
	viper.BindEnv("POSTGRES_SSL_MODE")
	viper.BindEnv("POSTGRES_TZ")

	// Binding SMTP
	viper.BindEnv("SMTP_HOST")
	viper.BindEnv("SMTP_HOST_USER")
	viper.BindEnv("SMTP_HOST_PASSWORD")
	viper.BindEnv("SMTP_PORT")
	viper.BindEnv("CS_EMAIL_ADDRESS")
}
