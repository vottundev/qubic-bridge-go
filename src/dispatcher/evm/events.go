package evm

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

var (
	waitGroup        sync.WaitGroup
	waitGroupStarted bool
)

func Susbscribe(client *ethclient.Client, contract *common.Address, info chan types.Log, errorChan chan<- error) {

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		subscribeToEvents(client, contract, info, errorChan)
	}()

	if !waitGroupStarted {
		go func() {
			waitGroupStarted = true
			waitGroup.Wait()
		}()
	}
}

func subscribeToEvents(client *ethclient.Client, contract *common.Address, info chan<- types.Log, errChan chan<- error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{*contract},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Errorf("Failed creating subscription: %+v", err)
		errChan <- err
		waitGroup.Done()
		return
	}

	for {
		select {
		case err := <-sub.Err():
			errChan <- err
		case vLog := <-logs:
			info <- vLog

		}
	}
}
