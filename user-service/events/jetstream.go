package events

import "github.com/nats-io/nats.go"

func CreateJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "USER_EVENTS",
		Subjects: []string{"user.created"},
	})

	if err != nil {
		return nil, err
	}
	return js, nil
}
