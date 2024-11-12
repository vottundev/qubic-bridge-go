package controller

import (
	"net/http"

	"github.com/vottundev/vottun-qubic-bridge-go/dto"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

const (
	ERROR_DECODING_DATA_FROM_REQUEST string = "failed decoding data from request"
)

func ProcessOrder(w http.ResponseWriter, r *http.Request) {

	order := &dto.OrderReceivedDTO{}

	err := decodeRequestIntoStruct(w, r, &order)

	if err != nil {
		log.Errorf("failed decoding request payload into struct")
		return
	}

	log.Tracef("order: %+v\n", order)

}
