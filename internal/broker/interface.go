package broker

type Broker interface {
	Publish(subject string, msg []byte) error
	Subscribe(subject string, handler func(msg []byte)) error
}
