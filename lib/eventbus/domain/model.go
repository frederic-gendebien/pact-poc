package domain

type EventDefinition interface {
	GetDomain() string
	GetName() string
	GetType() interface{}
}

type Event interface {
	GetDefinition() EventDefinition
	GetEntityId() string
	GetPayload() interface{}
}

type EventHandler interface {
	GetEventDefinition() EventDefinition
	ProcessEvent(interface{}) error
	HandleError(interface{}, error)
}
