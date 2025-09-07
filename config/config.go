package config

import (
	"strings"

	"log"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/mcuadros/go-defaults"
)

var C Config

func init() {
	k := koanf.New(".")
	defaults.SetDefaults(&C)

	if err := godotenv.Load(); err != nil {
		log.Println("error loading .env variables", err)
	}

	envProvider := env.Provider("", "__", strings.ToLower)
	err := k.Load(envProvider, nil)
	if err != nil {
		log.Fatal(err)
	}

	unmarshalerConfig := koanf.UnmarshalConf{Tag: "json"}
	if err := k.UnmarshalWithConf("", &C, unmarshalerConfig); err != nil {
		log.Fatal(err)
	}

	v := validator.New()
	if err := v.Struct(C); err != nil {
		log.Fatal(err)
	}
}
