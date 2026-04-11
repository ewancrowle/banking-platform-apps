package config

type Config struct {
	Port int `envconfig:"port" default:"8080"`
}
