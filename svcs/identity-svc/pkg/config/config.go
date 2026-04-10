package config

type Config struct {
	Port      int `envconfig:"port" default:"8080"`
	MachineID int `envconfig:"machine_id" required:"true"`
}
