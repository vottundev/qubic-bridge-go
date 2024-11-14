package cache

import (
	"context"

	"github.com/vottundev/vottun-qubic-bridge-go/config"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

type PubSubHandler func(channel string, payload string)

func subscribeToQubicEvens(ctx context.Context, pubSubHandler PubSubHandler) {
	if qubicRedisClient != nil {
		// Subscribe to the specified channel
		sub := qubicRedisClient.Subscribe(ctx, config.Config.Cache.QubicEventsChannel)

		// Ensure the subscription is closed when done
		defer sub.Close()

		// Get the channel to receive messages
		channel := sub.Channel()

		// Listen for messages in a loop until context is done
		for {
			select {
			case <-ctx.Done():
				// If the context is done, exit the loop
				log.Infoln("Context done, stopping subscription.")
				return
			case msg := <-channel:
				// Handle incoming messages
				pubSubHandler(msg.Channel, msg.Payload)
				log.Tracef("Received message: %s from channel: %s\n", msg.Payload, msg.Channel)
			}
		}
	}
}
