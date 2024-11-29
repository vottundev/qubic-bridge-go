package evm

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/vottundev/vottun-qubic-bridge-go/assets"
	"github.com/vottundev/vottun-qubic-bridge-go/config"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

const (
	ErrorVtnAlreadySubscribed    = "ERROR_VTN_ALREADY_SUBSCRIBED"
	ErrorVtnAlreadySubscribedMsg = "Already subscribed to $VTN events"
)

var (
	parsedAbi abi.ABI
	// subscribed bool
)

func SubscribeToEVMEvents(chainInfo *config.ChainInfo) {

	var err error

	parsedAbi, err = abi.JSON(bytes.NewReader(assets.TestEvent))
	if err != nil {
		log.Errorf("Failed parsing contract ABI EVM orders contract: %+v", err)
		return
		// return errors.New(constants.ErrorParsingAbi, fmt.Sprintf("Failed parsing contract ABI EVM orders contract%+v", err))
	}

	//1. get client to get gas price from the blockchain
	evmClient, err := getEthereumClient(chainInfo, true)
	if err != nil {
		log.Errorf("Error getting ethereum client. %+v", err)
		time.Sleep(5 * time.Second)

		go SubscribeToEVMEvents(chainInfo)
		return
		// return errors.New(constants.ErrorGettingEvmClient, fmt.Sprintf("Error getting ethereum client. %+v", err))
	}

	logs := make(chan types.Log)
	errChan := make(chan error)

	go Susbscribe(
		evmClient,
		chainInfo.ContractAddress,
		logs,
		errChan,
	)

	go func() {

		for {
			select {
			case err := <-errChan:
				log.Errorf("%+v ************************************************************************************************************************", err)
				go SubscribeToEVMEvents(chainInfo)
				break
			case info := <-logs:
				err := processLog(info)
				if err != nil {
					log.Errorf("%+v", err)
				}
			}

		}

	}()
	// subscribed = true

}

func processLog(info types.Log) error {

	decodedLogs := DecodeEventLogs([]*types.Log{&info}, &parsedAbi)

	b, _ := json.MarshalIndent(decodedLogs, "", "  ")
	log.Tracef("%s", string(b))
	log.Tracef("%+v", decodedLogs)

	// if err := datalayer.PersistVtnEvent(&info, b); err != nil {
	// 	log.Errorf("%+v", err)
	// 	return err
	// }

	return nil
}

func getEthereumClient(client *config.ChainInfo, wss bool) (*ethclient.Client, error) {

	var e *ethclient.Client
	var err error

	if wss {
		e, err = ethclient.Dial(client.WssUrl)
	} else {
		e, err = ethclient.Dial(client.RpcUrl)
	}
	if err != nil {
		log.Errorf("There is an error getting connection for network {%+v}: %+v", client, err)
		return nil, err
	}

	return e, nil
}
