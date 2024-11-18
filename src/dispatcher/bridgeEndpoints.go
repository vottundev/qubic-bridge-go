package dispatcher

import (
	"encoding/json"
	"net/http"

	"github.com/vottundev/vottun-qubic-bridge-go/config"
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
		// DispatchOrderForProcessing(order)
		grpc.ProcessQubicOrder(order)
	case dto.CONFIRM_ORDER:
	}
}

func DispatchOrderForProcessing(order *dto.OrderReceivedDTO) error {

	err := sendBridgeRequest(
		&RequestModel{
			Url:          config.Config.InternalEndpoints.ProcessOrder,
			HttpMethod:   http.MethodPost,
			RequestDto:   &order,
			ParseRequest: true,
		},
	)

	if err != nil {
		log.Errorf("An error has raised calling core api Create New Custodied Wallet. %+v", err)
		return err
	}

	return nil
}
