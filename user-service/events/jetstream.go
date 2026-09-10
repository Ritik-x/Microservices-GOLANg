package events

import "github.com/nats-io/nats.go"

func CreateJetStream(nc *nats.Conn) (nats.JetStreamContext, error) {

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// Check if stream already exists
	_, err = js.StreamInfo("USER_EVENTS")

	if err == nil {
		// Stream already exists, so use it
		return js, nil
	}

	// Stream does not exist, create it
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "USER_EVENTS",
		Subjects: []string{"user.created"},
	})

	if err != nil {
		return nil, err
	}

	return js, nil
}
