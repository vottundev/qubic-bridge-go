package evm

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

type AbiMethod struct {
	Inputs []struct {
		Indexed      bool   `json:"indexed"`
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"inputs,omitempty"`
	Outputs []struct {
		InternalType string `json:"internalType"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"outputs,omitempty"`
	Name            string `json:"name"`
	Anonymous       bool   `json:"anonymous,omitempty"`
	Type            string `json:"type"`
	StateMutability string `json:"stateMutability,omitempty"`
	Signature       string `json:"signature"`
	Keccak          string `json:"keccak"`
}

type AbiContract struct {
	Constructor AbiMethod
	Events      map[string]AbiMethod
	Methods     map[string]AbiMethod
	Errors      map[string]AbiMethod
}

func (a *AbiContract) LoadAbiFromSpecs(abiBytes []byte) error {
	var methods *[]AbiMethod

	err := json.Unmarshal(abiBytes, &methods)
	if err != nil {
		log.Errorf("%+v", err)
		return err
	}
	// a = &AbiContract{ContractSpecs: c}
	a.Events = make(map[string]AbiMethod)
	a.Methods = make(map[string]AbiMethod)
	a.Errors = make(map[string]AbiMethod)

	for _, m := range *methods {

		switch m.Type {
		case "constructor":
			a.Constructor = m
		case "event":
			a.Events[m.Name] = m
		case "function":
			a.Methods[m.Name] = m
		case "error":
			a.Errors[m.Name] = m

		default:
			log.Debugf("m.Type: %v\n", m.Type)
		}
	}

	return nil
}

func (a *AbiContract) ProcessMethodsForSubscription() {

	var mu = sync.RWMutex{}

	mu.Lock()
	defer mu.Unlock()

	for k, v := range a.Methods {
		v.Signature = v.GetFunctionSignature()
		v.Keccak = v.GetFunctionKeccak(0)
		a.Methods[k] = v
	}
}
func (a *AbiMethod) GetFunctionSignature() string {

	var params string
	r := strings.Builder{}

	// r.WriteString(a.Name)
	// r.WriteString("(")
	for _, p := range a.Inputs {
		r.WriteString(p.Type)
		r.WriteString(",")
	}

	if r.Len() > 1 {
		params = r.String()[0 : r.Len()-1]
	}
	return fmt.Sprintf(
		"%s(%s)",
		a.Name,
		params,
	)
}

// Returns the first 4 bytes of the function keccak with 0x prefix
func (a *AbiMethod) GetFunctionCode() string {
	return a.GetFunctionKeccak(4)
}

// Returns the "length" first bytes of the function keccak with 0x prefix.
//
// If length is set to zero, the the full keccak string is returned
//
//	params
//	  length - the number of bytes to return
func (a *AbiMethod) GetFunctionKeccak(length int) string {

	kc := crypto.Keccak256([]byte(a.GetFunctionSignature()))

	if length == 0 {
		length = len(kc)
	}

	return fmt.Sprintf("0x%s", hex.EncodeToString(kc[:length]))
}
