package config

import (
	"github.com/frederic-gendebien/pact-poc/lib/config/environment"
	"io"
	"log"
	"os"
)

const (
	Mode            = "CONFIGURATION_MODE"
	ModeEnvironment = "environment"
)

func NewConfiguration() Configuration {
	mode := os.Getenv("CONFIGURATION_MODE")
	switch mode {
	case "":
		log.Println("default configuration mode")
		fallthrough
	case ModeEnvironment:
		return environment.NewConfiguration()
	default:
		log.Fatalf("unknown configuration mode: %s", mode)
		return nil
	}
}

type Configuration interface {
	io.Closer
	GetString(name string, defaultProvider func() string) string
	GetStringOrCrash(name string) string
}
