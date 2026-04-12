package config

type Config struct {
	Port int `envconfig:"port" default:"8080"`

	IdentityServiceAddr string `envconfig:"identity_service_addr" required:"true"`

	DBHost     string `envconfig:"db_host" required:"true"`
	DBName     string `envconfig:"db_name" required:"true"`
	DBUsername string `envconfig:"db_username" required:"true"`
	DBPassword string `envconfig:"db_password" required:"true"`
}
