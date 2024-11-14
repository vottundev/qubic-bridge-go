package dispatcher

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/vottundev/vottun-qubic-bridge-go/config"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/net"
)

const (
	CONTENT_TYPE   string = "Content-Type"
	AUTH_APP_ID    string = "x-application-vkn"
	AUTHORIZATION  string = "Authorization"
	MIME_TYPE_JSON string = "application/json; charset=UTF-8"
)

type RequestModel struct {
	Url            string
	HttpMethod     string
	RequestDto     interface{}
	ResponseDto    interface{}
	TokenAuth      *string
	AppID          *string
	ResponseStatus int
	ParseRequest   bool
	ParseResponse  bool
}

func sendBridgeRequest(s *RequestModel) error {

	kk := &net.RequestApiEndpointInfo{
		EndpointUrl:  config.Config.InternalEndpoints.Host + s.Url,
		RequestData:  s.RequestDto,
		ResponseData: s.ResponseDto,
		HttpMethod:   s.HttpMethod,
		TokenAuth:    s.TokenAuth,
		AppID:        s.AppID,
	}
	err := net.RequestApiEndpoint(
		kk,
		setReqHeaders,
		s.ParseRequest,
		s.ParseResponse,
	)

	if err != nil {
		log.Errorf("An error has raised calling internal bridge endpoint. %+v", err)
		log.Errorf("%+v", *s)
		return err
	}

	s.ResponseDto = kk.ResponseData
	s.ResponseStatus = kk.ResponseStatus
	return err
}

func setReqHeaders(req *http.Request, tokenAuth, appID *string) {

	req.Header.Add(CONTENT_TYPE, MIME_TYPE_JSON)

	if tokenAuth != nil {
		req.Header.Add(AUTHORIZATION, fmt.Sprintf("Bearer %s", *tokenAuth))
	}
	if appID != nil {
		req.Header.Add(AUTH_APP_ID, *appID)

	}

}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
