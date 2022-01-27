package inmemory

import (
	"github.com/frederic-gendebien/pact-poc/lib/eventbus/domain"
	"math/rand"
)

type HandlerGroup []domain.EventHandler

func (g HandlerGroup) RandomHandler() domain.EventHandler {
	return g[int(rand.Uint32())%len(g)]
}

func NewHandlerGroups() HandlerGroups {
	return make(map[string]HandlerGroup)
}

type HandlerGroups map[string]HandlerGroup

func (h HandlerGroups) AddEventHandler(name string, handler domain.EventHandler) {
	h[name] = append(h[name], handler)
}

func (h HandlerGroups) SelectHandlers() []domain.EventHandler {
	selectedHandlers := make([]domain.EventHandler, 0, 2)
	for _, handlerGroup := range h {
		selectedHandlers = append(selectedHandlers, handlerGroup.RandomHandler())
	}

	return selectedHandlers
}
