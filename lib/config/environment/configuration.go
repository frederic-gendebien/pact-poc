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

func (c Configuration) GetOptionalValue(name string, defaultProvider func() string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultProvider()
	}

	return value
}

func (c Configuration) GetMandatoryValue(name string) string {
	return c.GetOptionalValue(name, func() string {
		log.Fatalf("missing mandatory property: %s", name)
		return ""
	})
}
