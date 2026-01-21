package valkey

import (
	"context"
	"errors"

	glide "github.com/valkey-io/valkey-glide/go/v2"
	glideconfig "github.com/valkey-io/valkey-glide/go/v2/config"
)

func Connect() (*glide.Client, error) {
	host := "localhost"
	port := 6379

	config := glideconfig.NewClientConfiguration().WithAddress(
		&glideconfig.NodeAddress{
			Host: host,
			Port: port,
		},
	)

	client, err := glide.NewClient(config)
	if err != nil {
		return nil, err
	}

	res, err := client.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	if res != "PONG" {
		return nil, errors.New("res not equal PONG")
	}

	return client, nil
}
