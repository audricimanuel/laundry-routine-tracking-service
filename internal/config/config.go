package config

type (
	Config struct {
		Env                   string     `mapstructure:"ENV"`
		SwaggerUsername       string     `mapstructure:"SWAGGER_USERNAME"`
		SwaggerPassword       string     `mapstructure:"SWAGGER_PASSWORD"`
		JWTSecret             string     `mapstructure:"JWT_SECRET"`
		JWTExpirationDuration float64    `mapstructure:"JWT_EXPIRATION_DURATION"`
		Host                  Host       `mapstructure:",squash"`
		DataSource            DataSource `mapstructure:",squash"`
		SMTPConfig            SMTPConfig `mapstructure:",squash"`
	}

	Host struct {
		Address      string `mapstructure:"HOST_ADDRESS"`
		Port         string `mapstructure:"HOST_PORT"`
		WriteTimeout int    `mapstructure:"HOST_WRITE_TIMEOUT"`
		ReadTimeout  int    `mapstructure:"HOST_READ_TIMEOUT"`
		IdleTimeout  int    `mapstructure:"HOST_IDLE_TIMEOUT"`
		FEBaseUrl    string `mapstructure:"FE_BASE_URL"`
	}

	DataSource struct {
		PostgresDBConfig PostgresDBConfig `mapstructure:",squash"`
	}

	PostgresDBConfig struct {
		Host     string `mapstructure:"POSTGRES_DB_HOST"`
		User     string `mapstructure:"POSTGRES_DB_USER"`
		Password string `mapstructure:"POSTGRES_DB_PASSWORD"`
		Name     string `mapstructure:"POSTGRES_DB_NAME"`
		Port     string `mapstructure:"POSTGRES_DB_PORT"`
		SSLMode  string `mapstructure:"POSTGRES_SSL_MODE"`
		Timezone string `mapstructure:"POSTGRES_TZ"`
	}

	SMTPConfig struct {
		Host           string `mapstructure:"SMTP_HOST"`
		User           string `mapstructure:"SMTP_HOST_USER"`
		Password       string `mapstructure:"SMTP_HOST_PASSWORD"`
		Port           string `mapstructure:"SMTP_PORT"`
		CSEmailAddress string `mapstructure:"CS_EMAIL_ADDRESS"`
	}
)
