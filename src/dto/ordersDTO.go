package dto

type PubSubEvent string

// used to get event payload as byte array
type BytesJsonData []byte

const (
	NEW_ORDER     PubSubEvent = "order"
	CONFIRM_ORDER PubSubEvent = "confirm"
)

type OrderReceivedDTO struct {
	OrderID            uint64 `json:"orderId"`
	OriginAccount      string `json:"originAccount"`
	DestinationAccount string `json:"destinationAccount"`
	Amount             string `json:"amount"`
	// Memo               string `json:"memo"`
	SourceChain uint32 `json:"sourceChain"`
}

type RedisPubSubDTO struct {
	EventType PubSubEvent   `json:"eventType"`
	Payload   BytesJsonData `json:"payload"`
}

func (r *BytesJsonData) UnmarshalJSON(data []byte) error {
	*r = append((*r)[0:0], data...)

	return nil
}
