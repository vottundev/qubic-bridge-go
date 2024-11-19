package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

type LogTopic struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LogData struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
type TransactionLog struct {
	Address string     `json:"address"`
	Name    string     `json:"name"`
	Topics  []LogTopic `json:"topics"`
	Data    []LogData  `json:"data"`
}

// Decodes the event logs using ABI if exists
func DecodeEventLogs(logs []*types.Log, contractABI *abi.ABI) []TransactionLog {
	if contractABI == nil {
		return nil
	}
	return decodeLogEventWithAbi(logs, contractABI)
}

func decodeLogEventWithAbi(logs []*types.Log, contractABI *abi.ABI) []TransactionLog {

	// We add a recover function from panics to prevent our API from crashing due to an unexpected error
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()

	result := make([]TransactionLog, 0)

	for _, tLog := range logs {

		decoded, err := contractABI.EventByID(tLog.Topics[0])

		if err != nil {
			log.Errorf("Error decoding log: %+v", err)
			continue
		}
		decodedData, err := decoded.Inputs.Unpack(tLog.Data)
		if err != nil {
			log.Errorf("Error unpacking inputs: %+v", err)
			continue
		}

		t := TransactionLog{}
		t.Address = tLog.Address.Hex()

		t.Name = decoded.RawName
		t.Topics = make([]LogTopic, 0)
		t.Data = make([]LogData, 0)

		var indexed, nonIndexed int = 1, 0
		for _, n := range decoded.Inputs {
			if n.Indexed {
				topic := LogTopic{
					Name: n.Name,
				}
				if n.Type.T == abi.AddressTy {
					topic.Value = common.HexToAddress(tLog.Topics[indexed].Hex()).Hex()
				} else {
					topic.Value = tLog.Topics[indexed].Hex()
				}
				indexed++
				t.Topics = append(t.Topics, topic)

			} else {
				data := LogData{
					Name:  n.Name,
					Value: decodedData[nonIndexed],
				}
				nonIndexed++
				t.Data = append(t.Data, data)
			}
		}

		result = append(result, t)
	}

	return result
}
