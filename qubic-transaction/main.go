package main

import (
	"bytes"
	"encoding/hex"
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
	err := run()
	if err != nil {
		log.Fatalf("got err: %s", err.Error())
	}
}

func run() error {
	sourceID := "PXABYVDPJRRDKELEYSHZWJCBEFJCNERNKKUWXHANCDPQEFGDIUGUGAUBBCYK"      // my ID (current adminID)
	destinationID := "LAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKPTJ" // contract ID
	adminID := "FPOBXDXMCSNZLCFWIERCYWYSQBLDPXRMHSDKVIFMTESIUGNFLAZLDLTDYMGF"       // New admin ID

	srcID := types.Identity(sourceID)
	destID := types.Identity(destinationID)
	admID := types.Identity(adminID)
	srcPubKey, err := srcID.ToPubKey(false)
	if err != nil {
		return errors.Wrap(err, "converting src id string to pubkey")
	}
	destPubKey, err := destID.ToPubKey(false)
	if err != nil {
		return errors.Wrap(err, "converting dest id string to pubkey")
	}
	adminPubKey, err := admID.ToPubKey(false)
	if err != nil {
		return errors.Wrap(err, "converting admin id string to pubkey")
	}

	tx := types.Transaction{
		SourcePublicKey:      srcPubKey,
		DestinationPublicKey: destPubKey,
		Amount:               0,
		Tick:                 17420050, // this should be set to currentTick of node + 10
		InputType:            3,        //setAdmin input type 3
		InputSize:            32,       // admin address pubkey (32 bytes)
		Input:                adminPubKey[:],
	}
	fmt.Printf("source pubkey: %s\n", hex.EncodeToString(tx.SourcePublicKey[:]))

	unsignedDigest, err := tx.GetUnsignedDigest()
	if err != nil {
		log.Fatalf("got err: %s when getting unsigned digest local", err.Error())
	}

	subSeed, err := types.GetSubSeed("hnwizobbscgmqfvosckyrrigrvtfkucaqagguxscmyewagynkrwychd")
	if err != nil {
		log.Fatalf("got err %s when getting subSeed", err.Error())
	}

	sig, err := schnorrq.Sign(subSeed, tx.SourcePublicKey, unsignedDigest)
	if err != nil {
		log.Fatalf("got err: %s when signing", err.Error())
	}
	fmt.Printf("sig: %s\n", hex.EncodeToString(sig[:]))
	tx.Signature = sig

	encodedTx, err := tx.EncodeToBase64()
	if err != nil {
		log.Fatalf("got err: %s when encoding tx to base 64", err.Error())
	}
	fmt.Printf("encodedTx: %s\n", encodedTx)

	id, err := tx.ID()
	if err != nil {
		log.Fatalf("got err: %s when getting tx id", err.Error())
	}

	fmt.Printf("tx id(hash): %s\n", id)

	url := baseUrl + "/broadcast-transaction"
	payload := struct {
		EncodedTransaction string `json:"encodedTransaction"`
	}{
		EncodedTransaction: encodedTx,
	}

	buff := new(bytes.Buffer)
	err = json.NewEncoder(buff).Encode(payload)
	if err != nil {
		log.Fatalf("got err: %s when encoding payload", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, url, buff)
	if err != nil {
		log.Fatalf("got err: %s when creating request", err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("got err: %s when performing request", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Got non 200 status code: %d", res.StatusCode)
	}

	type response struct {
		PeersBroadcasted   uint32 `json:"peersBroadcasted"`
		EncodedTransaction string `json:"encodedTransaction"`
		TxID               string `json:"transactionId"`
	}
	var body response

	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		log.Fatalf("got err: %s when decoding body", err.Error())
	}

	fmt.Printf("%+v\n", body)

	return nil
}
