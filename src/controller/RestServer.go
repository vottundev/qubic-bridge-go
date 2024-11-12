package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/vottundev/vottun-qubic-bridge-go/config"
	"github.com/vottundev/vottun-qubic-bridge-go/controller/interceptor"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/net"
)

var (
	httpRouter *mux.Router
	httpPort   string
	httpServer *http.Server
)

// MiddlewareHandlerFunc builds on top of http.HandlerFunc, and exposes API to intercept with MiddlewareInterceptor.
// This allows building complex long chains without complicated struct manipulation
type MiddlewareHandlerFunc http.HandlerFunc

func prepareCors(router *mux.Router) http.Handler {

	c := cors.New(cors.Options{
		AllowedOrigins: config.Config.Cors.AllowedOrigins,
		AllowedMethods: config.Config.Cors.AllowedMethods,
		AllowedHeaders: config.Config.Cors.AllowedHeaders,
		ExposedHeaders: []string{},
	})

	return c.Handler(router)
}

func SetupRestServer(port string) {

	var err error

	httpRouter = mux.NewRouter()
	httpPort = port

	httpRouter.Path(config.Config.Http.Route + "/ping").Methods(http.MethodGet).HandlerFunc(http.HandlerFunc(MiddlewareHandlerFunc(IsAlive)))

	log.Infof("Vottun Qubic Public service Listening on port %s", port)

	httpRouter.Use(interceptor.NewElapsedTimeInterceptor())

	if log.LogLevel <= log.DEBUG {
		httpServer, err = net.ListenAndServe(
			net.ListenAndServeInfo{
				Ipversion: net.IPV4,
				Address:   "0.0.0.0:" + port,
				Handler:   prepareCors(httpRouter),
			})
	} else {
		httpServer, err = net.ListenAndServe(
			net.ListenAndServeInfo{
				Ipversion:    net.IPV4,
				Address:      "0.0.0.0:" + port,
				Handler:      prepareCors(httpRouter),
				WriteTimeout: 30,
				ReadTimeout:  15,
				IdleTimeout:  60,
			})
	}

	if err != nil {
		log.Errorf("failed starting public http server: %+v", err)
		panic(err)
	}
}
func RestartServer() {

	ShutdownHttpServer("Requested public http server restart from admin")
	SetupRestServer(httpPort)
}

func ShutdownHttpServer(reason string) {
	net.ShutDown(httpServer, reason)
}
