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
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

type Arguments struct {
	Secret    string            `conf:"env:SECRET,required"`
	YamlFile  string            `conf:"flag:yaml,short:y,required"`
	LogStdout bool              `conf:"flag:logstdout,short:s"`
	LogFile   *string           `conf:"flag:logfile,short:f"`
	Port      string            `conf:"flag:port,short:p"`
	LogLevel  log.LogLevelValue `conf:"default:INFO,flag:loglevel,short:l"`
}

var (
	args    *Arguments
	sigterm chan bool
)

func main() {
	log.Infoln("Starting Cache")
	cache.Start()

	go controller.SetupRestServer(args.Port)

	log.Infoln("Telegram Vottun Dojo Up'n'Ready.")

	_, cancel := context.WithCancel(context.Background())

	handleSigTerm()

	receivedSigterm := <-sigterm
	if receivedSigterm {
		log.Infoln("Sigterm arrived. Shut down")
		controller.ShutdownHttpServer("SigTerm Received")
		cache.StopRedisClients()
		cancel()
	}
}

func init() {

	config.ExecutionTime = time.Now()

	args = &Arguments{}

	if err := conf.Parse(os.Args[1:], "BRIDGE", &args); err != nil {
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
