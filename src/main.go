package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/vottundev/vottun-qubic-bridge-go/cache"
	"github.com/vottundev/vottun-qubic-bridge-go/config"
	"github.com/vottundev/vottun-qubic-bridge-go/controller"
	"github.com/vottundev/vottun-qubic-bridge-go/dispatcher"
	"github.com/vottundev/vottun-qubic-bridge-go/dispatcher/evm"
	"github.com/vottundev/vottun-qubic-bridge-go/dto"
	"github.com/vottundev/vottun-qubic-bridge-go/grpc"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

type ExecutionType string

const (
	EXECUTION_TYPE_BRIDGE     ExecutionType = "bridge"
	EXECUTION_TYPE_DISPATCHER ExecutionType = "dispatcher"
)

type Arguments struct {
	Secret         string            `conf:"env:SECRET,required"`
	YamlFile       string            `conf:"flag:yaml,short:y,required"`
	LogStdout      bool              `conf:"flag:logstdout,short:s"`
	LogFile        *string           `conf:"flag:logfile,short:f"`
	Port           string            `conf:"flag:port,short:p"`
	InternalPort   string            `conf:"flag:internalport,short:i"`
	GrpcServerPort uint16            `conf:"default:50551,flag:grpc-server-port"`
	LogLevel       log.LogLevelValue `conf:"default:INFO,flag:loglevel,short:l"`
}

type ExecutionArgument struct {
	Execution ExecutionType `conf:"flag:launch,short:u,default:bridge"`
}

var (
	args         *Arguments
	execArgument *ExecutionArgument
	sigterm      chan bool
	cancel       context.CancelFunc
)

func main() {
	fmt.Println("Starting")
	config.ExecutionTime = time.Now()

	execArgument = &ExecutionArgument{}

	//check for execution argument to start app as bridge or redis dispatcher
	if err := conf.Parse(os.Args[1:], "BRIDGE", execArgument); err != nil {
		log.Panic(err)
	}

	_, cancel = context.WithCancel(context.Background())

	switch execArgument.Execution {
	case EXECUTION_TYPE_BRIDGE:
		mainBridge()
	case EXECUTION_TYPE_DISPATCHER:
		mainDispatcher()
	}
	handleSigTerm()

	order := dto.OrderReceivedDTO{
		OrderID:            1,
		OriginAccount:      "AAAAAAAAAAAAAAAAAAAAAA",
		DestinationAccount: "0x123412341341",
		Amount:             "3456789",
	}

	p := map[string]interface{}{
		"eventType": dto.NEW_ORDER,
		"payload":   order,
	}
	b, _ := json.Marshal(p)

	fmt.Printf("%s", string(b))

	receivedSigterm := <-sigterm
	if receivedSigterm {
		log.Infoln("Sigterm arrived. Shut down")
		if execArgument.Execution == EXECUTION_TYPE_BRIDGE {
			controller.ShutdownHttpServer("SigTerm Received")
			grpc.StopGrpcServer()
		} else if execArgument.Execution == EXECUTION_TYPE_DISPATCHER {
			grpc.StopGrprClientConnection()
		}
		cache.StopRedisClients()
		cancel()
	}
}

func mainDispatcher() {
	log.Infof("Begin service as Dispatcher")
	parseBridgeArguments()

	log.Infoln("Starting Cache")
	cache.Start(false, dispatcher.PubSubHandler)
	go grpc.StartGrpcClientConnection(args.GrpcServerPort)
	evm.SubscribeToEVMEvents(config.Config.Evm.Chains[config.CHAIN_ARB])

}
func mainBridge() {
	log.Infof("Begin service as Bridge")
	parseBridgeArguments()

	log.Infoln("Starting Cache")
	cache.Start(true, nil)

	go controller.SetupRestServer(args.Port)
	go grpc.StartGrpcServer(args.GrpcServerPort)

	log.Infoln("Telegram Vottun Dojo Up'n'Ready.")

}

func parseBridgeArguments() {
	args = &Arguments{}

	if err := conf.Parse(os.Args[1:], "BRIDGE", args); err != nil {
		log.Panic(err)
	}

	config.CreateProperties(args.YamlFile, args.Secret)

	setLog(args.LogFile, args.LogStdout, args.LogLevel)
	if log.LogLevel == log.TRACE {
		b, _ := json.MarshalIndent(config.Config, "", "  ")
		log.Tracef("Properties:\n%s", string(b))
	}
	log.Infoln("Properties, Environment and Log config processed.")
}

func setLog(logFile *string, stdout bool, logLevel log.LogLevelValue) {
	if stdout {
		log.SetOutput(os.Stdout)
	} else {
		// If the file doesn't exist, create it or append to the file
		file, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
	}
	log.LogLevel = logLevel
}

func handleSigTerm() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sigterm = make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		sigterm <- true
	}()
}
