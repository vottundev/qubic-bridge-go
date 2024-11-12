package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vottundev/vottun-qubic-bridge-go/config"
	"github.com/vottundev/vottun-qubic-bridge-go/controller/interceptor"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/net"
)

var (
	httpInternalRouter *mux.Router
	httpInternalPort   string
	httpInternalServer *http.Server
)

func SetupInternalRestServer(port string) {

	var err error

	httpInternalRouter = mux.NewRouter()
	httpInternalPort = port

	httpInternalRouter.Path(config.Config.Http.InternalRoute + "/ping").Methods(http.MethodGet).HandlerFunc(http.HandlerFunc(MiddlewareHandlerFunc(IsAlive)))
	httpInternalRouter.Path(config.Config.Http.InternalRoute + "/order").Methods(http.MethodPost).HandlerFunc(http.HandlerFunc(MiddlewareHandlerFunc(ProcessOrder)))

	log.Infof("Vottun Qubic Private service Listening on port %s", port)

	httpInternalRouter.Use(interceptor.NewElapsedTimeInterceptor())

	if log.LogLevel <= log.DEBUG {
		httpInternalServer, err = net.ListenAndServe(
			net.ListenAndServeInfo{
				Ipversion: net.IPV4,
				Address:   "0.0.0.0:" + port,
				Handler:   httpInternalRouter,
			})
	} else {
		httpInternalServer, err = net.ListenAndServe(
			net.ListenAndServeInfo{
				Ipversion:    net.IPV4,
				Address:      "0.0.0.0:" + port,
				WriteTimeout: 30,
				ReadTimeout:  15,
				IdleTimeout:  60,
				Handler:      httpInternalRouter,
			})
	}

	if err != nil {
		log.Errorf("failed starting internal http server: %+v", err)
		panic(err)
	}
}

func RestartInternalServer() {

	ShutdownInternalHttpServer("Requested internal http server restart from admin")
	SetupInternalRestServer(httpInternalPort)
}

func ShutdownInternalHttpServer(reason string) {
	net.ShutDown(httpInternalServer, reason)
}
