package config

type Config struct {
	Port int `envconfig:"port" default:"8080"`

	PaymentServiceAddr string `envconfig:"payment_service_addr" required:"true"`
}
