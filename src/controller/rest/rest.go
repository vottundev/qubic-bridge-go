package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-http-utils/headers"
)

type ErrorDTO struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// NewErrorDTO returns a new ErrorDTO
func NewErrorDTO(code, message string) ErrorDTO {
	e := ErrorDTO{code, message, time.Now()}
	return e
}

// Returns a Raw error writting the code and the message received
func ReturnError(w http.ResponseWriter, code, message string, status int) {

	w.Header().Add(headers.ContentType, "application/json")

	errorDTO := NewErrorDTO(code, message)
	w.WriteHeader(status)
	errorMessage, _ := json.Marshal(errorDTO)

	fmt.Fprint(w, string(errorMessage))

}

func ReturnResponseToClient(w http.ResponseWriter, value interface{}) {
	ReturnResponseToClientWithStatus(w, value, http.StatusOK)
}

func ReturnResponseToClientWithStatus(w http.ResponseWriter, value interface{}, httpStatus int) {
	w.Header().Add(headers.ContentType, "application/json")
	w.WriteHeader(httpStatus)
	b, err := json.Marshal(value)
	if err != nil {
		//TODO: Error marshalling
	} else {
		fmt.Fprint(w, string(b))
	}

}
