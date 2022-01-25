package domain

type EventDefinition interface {
	GetDomain() string
	GetName() string
}

type Event interface {
	GetDefinition() EventDefinition
	GetEntityId() string
	GetPayload() interface{}
}

type EventHandler interface {
	GetName() string
	GetEvent() EventDefinition
	GetHandling() Handler
	GetErrorHandling() ErrorHandler
}

type Handler func(interface{}) error
type ErrorHandler func(interface{}, error)
