package controller

import (
	"net/http"
	"runtime"

	"github.com/vottundev/vottun-qubic-bridge-go/config"
	"github.com/vottundev/vottun-qubic-bridge-go/controller/rest"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

func IsAlive(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			log.Errorf("There has been an unmanaged error at IsAlive: %+v", err)
		}
	}()

	mm := make(map[string]interface{})
	p := make(map[string]interface{})

	mm["executionTime"] = config.ExecutionTime

	p["APP-NAME"] = "VOTTUN-QUBIC-BRIDGE"
	p["APP-VERSION"] = "v1"
	p["GOARCH"] = runtime.GOARCH
	p["GOOS"] = runtime.GOOS
	p["GOROOT"] = runtime.GOROOT()
	p["NumCPU"] = runtime.NumCPU()
	p["NumGoroutine"] = runtime.NumGoroutine()
	p["Version"] = runtime.Version()

	mm["runtime"] = p
	rest.ReturnResponseToClient(w, mm)
}
