package config

type Config struct {
	Port               int    `envconfig:"port" default:"8080"`
	AccountServiceAddr string `envconfig:"account_service_addr" required:"true"`
	ClickHouseURL      string `envconfig:"clickhouse_url" required:"true"`
}
