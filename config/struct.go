package config

type Config struct {
	Postgres Postgres `json:"postgres"`
	JWT      JWT      `json:"jwt"`
}

type Postgres struct {
	Host     string `json:"host" default:"localhost"`
	Port     string `json:"port" default:"5432"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	DB       string `json:"db" default:"noroask"`
	Timezone string `json:"time_zone" default:"UTC"`
	SSLMode  string `json:"ssl_mode" default:"disable"`
}

type JWT struct {
	Secret 	string `json:"secret" validate:"required"`
}