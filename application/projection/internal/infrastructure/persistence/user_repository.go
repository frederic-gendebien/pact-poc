package persistence

import (
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/domain/repository"
	"github.com/frederic-gendebien/pact-poc/application/projection/internal/infrastructure/persistence/inmemory"
	"github.com/frederic-gendebien/pact-poc/lib/config"
	"log"
)

const (
	Mode         = "PERSISTENCE_MODE"
	ModeInMemory = "inmemory"
)

func NewUserRepository(configuration config.Configuration) repository.UserRepository {
	mode := configuration.GetStringOrCrash(Mode)
	switch mode {
	case ModeInMemory:
		return inmemory.NewUserRepository()
	default:
		log.Fatalf("unknown persistence mode: %s", mode)
		return nil
	}
}
