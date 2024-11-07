package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/vottundev/vottun-qubic-bridge-go/config"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/net"
)

var (
	httpRouter *mux.Router
	httpPort   string
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

	httpRouter = mux.NewRouter()
	httpPort = port

	httpRouter.Path(config.Config.Http.Route + "/ping").Methods(http.MethodGet).HandlerFunc(http.HandlerFunc(MiddlewareHandlerFunc(IsAlive)))

	log.Infof("Telegram Vottun Dojo service Listening on port %s", port)

	if log.LogLevel <= log.DEBUG {
		log.Fatal(net.ListenAndServe(
			net.ListenAndServeInfo{
				Ipversion: net.IPV4,
				Address:   "0.0.0.0:" + port,
				Handler:   prepareCors(httpRouter),
			}))
	} else {
		log.Fatal(net.ListenAndServe(
			net.ListenAndServeInfo{
				Ipversion:    net.IPV4,
				Address:      "0.0.0.0:" + port,
				Handler:      prepareCors(httpRouter),
				WriteTimeout: 15,
				ReadTimeout:  15,
				IdleTimeout:  60,
			}))
	}
}
func RestartServer() {

	ShutdownHttpServer("Requested restart from admin")
	SetupRestServer(httpPort)
}

func ShutdownHttpServer(reason string) {
	net.ShutDown(reason)
}
