package dispatcher

import (
	"encoding/json"

	"github.com/vottundev/vottun-qubic-bridge-go/dto"
	"github.com/vottundev/vottun-qubic-bridge-go/grpc"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

func PubSubHandler(channel string, payload string) {

	msg := &dto.RedisPubSubDTO{}

	err := json.Unmarshal([]byte(payload), &msg)
	if err != nil {
		log.Errorf("failed unmarshaling message payload: %+v", err)
		return
	}

	switch msg.EventType {
	case dto.NEW_ORDER:
		order := &dto.OrderReceivedDTO{}
		err = json.Unmarshal(msg.Payload, &order)
		if err != nil {
			log.Errorf("failed unmarshaling order: %+v", err)
			return
		}

		go grpc.ProcessQubicOrder(order)
	case dto.CONFIRM_ORDER:
	}
}
