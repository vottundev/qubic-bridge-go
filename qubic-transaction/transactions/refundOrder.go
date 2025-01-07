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
	err := refundOrder(0)
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

func refundOrder(orderID uint64) error {
	sourceID := "PXABYVDPJRRDKELEYSHZWJCBEFJCNERNKKUWXHANCDPQEFGDIUGUGAUBBCYK"
	destinationID := "LAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKPTJ"

	srcID := types.Identity(sourceID)
	destID := types.Identity(destinationID)

	srcPubKey, err := srcID.ToPubKey(false)
	if err != nil {
		return errors.Wrap(err, "converting src id to pubkey")
	}
	destPubKey, err := destID.ToPubKey(false)
	if err != nil {
		return errors.Wrap(err, "converting dest id to pubkey")
	}

	orderBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(orderBytes, orderID) // 8 bytes para uint64

	tx := types.Transaction{
		SourcePublicKey:      srcPubKey,
		DestinationPublicKey: destPubKey,
		Amount:               0,
		Tick:                 17420070,
		InputType:            7, // refundOrder
		InputSize:            8,
		Input:                orderBytes,
	}

	return signAndBroadcast(tx)
}
