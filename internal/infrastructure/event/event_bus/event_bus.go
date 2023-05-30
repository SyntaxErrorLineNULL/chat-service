package event_bus

type EventBus interface {
	Subscribe() error
	Publish(eventName string, args ...interface{}) any
	Unsubscribe(eventName string) error
}
