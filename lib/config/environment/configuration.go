package environment

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func NewConfiguration() Configuration {
	log.Println("starting inmemory configuration")
	_ = godotenv.Load()
	return Configuration{}
}

type Configuration struct {
}

func (c Configuration) Close() error {
	log.Println("closing inmemory configuration")

	return nil
}

func (c Configuration) GetString(name string, defaultProvider func() string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultProvider()
	}

	return value
}

func (c Configuration) GetStringOrCrash(name string) string {
	return c.GetString(name, func() string {
		log.Fatalf("missing mandatory property: %s", name)
		return ""
	})
}
