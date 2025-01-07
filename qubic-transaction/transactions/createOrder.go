package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/qubic/go-node-connector/types"
	"github.com/qubic/go-schnorrq"
)

const baseUrl = "http://127.0.0.1:80/v1"

func main() {
	err := createOrder("MNGERABCDXPLKWYSQTRFHDLKJVRFMNDQLPKFDIFNSRTQLWOSJXQPLKTZMQGJ", 2000, true)
	if err != nil {
		log.Fatalf("got err: %s", err.Error())
	}
}

func signAndBroadcast(tx types.Transaction) error {
	unsignedDigest, err := tx.GetUnsignedDigest()
	if err != nil {
		return errors.Wrap(err, "getting unsigned digest")
	}

	subSeed, err := types.GetSubSeed("hnwizobbscgmqfvosckyrrigrvtfkucaqagguxscmyewagynkrwychd")
	if err != nil {
		return errors.Wrap(err, "getting subSeed")
	}

	sig, err := schnorrq.Sign(subSeed, tx.SourcePublicKey, unsignedDigest)
	if err != nil {
		return errors.Wrap(err, "signing transaction")
	}
	tx.Signature = sig

	encodedTx, err := tx.EncodeToBase64()
	if err != nil {
		return errors.Wrap(err, "encoding transaction")
	}

	url := baseUrl + "/broadcast-transaction"
	payload := map[string]string{"encodedTransaction": encodedTx}

	buff := new(bytes.Buffer)
	json.NewEncoder(buff).Encode(payload)

	req, err := http.NewRequest(http.MethodPost, url, buff)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "performing request")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("broadcast failed with status: %d", res.StatusCode)
	}

	return nil
}

func createOrder(ethAddress string, amount uint64, fromQubicToEthereum bool) error {
	sourceID := "PXABYVDPJRRDKELEYSHZWJCBEFJCNERNKKUWXHANCDPQEFGDIUGUGAUBBCYK"
	destinationID := "LAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKPTJ"

	srcID := types.Identity(sourceID)
	destID := types.Identity(destinationID)
	ethID := types.Identity(ethAddress)

	srcPubKey, err := srcID.ToPubKey(false)
	if err != nil {
		return errors.Wrap(err, "converting src id to pubkey")
	}
	destPubKey, err := destID.ToPubKey(false)
	if err != nil {
		return errors.Wrap(err, "converting dest id to pubkey")
	}
	ethPubKey, err := ethID.ToPubKey(false)
	if err != nil {
		return errors.Wrap(err, "converting eth id to pubkey")
	}

	amountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(amountBytes, amount) // 8 bytes para uint64
	directionByte := []byte{0}
	if fromQubicToEthereum {
		directionByte[0] = 1
	}

	input := append(ethPubKey[:], append(amountBytes, directionByte...)...) // 49 + 8 + 1 bytes

	tx := types.Transaction{
		SourcePublicKey:      srcPubKey,
		DestinationPublicKey: destPubKey,
		Amount:               amount,
		Tick:                 17420090,
		InputType:            1, // createOrder
		InputSize:            uint16(len(input)),
		Input:                input,
	}

	return signAndBroadcast(tx)
}
