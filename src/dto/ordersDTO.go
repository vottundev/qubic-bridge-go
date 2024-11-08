package dto

type OrderReceivedDTO struct {
	OrderID            string `json:"orderId"`
	OriginChain        uint32 `json:"origiChain"`
	OriginAccount      string `json:"originAccount"`
	DestinationAccount string `json:"destinationAccount"`
	Amount             string `json:"amount"`
	Memo               string `json:"memo"`
}
