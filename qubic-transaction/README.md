### Procedure methods

 All procedure methods require to execute transactions. So, first of all, connect to the Qubic node using `qubic-http` docker, and check the 'tick-info': 
 `curl http://127.0.0.1:80/v1/tick-info`

 Then, modify the 'Tick' value of the Transaction configured:

 	tx := types.Transaction{
		SourcePublicKey:      srcPubKey,
		DestinationPublicKey: destPubKey,
		Amount:               0,
		Tick:                 17420050, // this should be set to currentTick of node + 10
		InputType:            4,        // addManager input type 4
		InputSize:            32,       // manager address pubkey (32 bytes)
		Input:                managerPubKey[:],
	}

You should also verify that sourceID (line 26) and subseed (line 62) correspond to the invocator account signing the transaction.

Procedures like completeOrder, refundOrder and transfer require to pay a a fee. Then, modify the field 'Amount' of Transaction when its configured. 

---

### Function methods
All function methods can be called with curl requests using `qubic-http` docker. 

## 2. getOrder (uint64 input)

curl -X 'POST' \
'http://127.0.0.1:80/v1/querySmartContract' \
-H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"contractIndex": 11,
"inputType": 2,
"inputSize": 8,
"requestData": "AAAAAAAAAAA="
}'


## 11. getTotalReceivedTokens (no input, uint64 output)

curl -X 'POST' \
'http://127.0.0.1:80/v1/querySmartContract' \
-H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"contractIndex": 11,
"inputType": 11,
"inputSize": 0,
"requestData": ""
}'

## 12. getAdmin (no input, id(32 bytes) output)

curl -X 'POST' \
'http://127.0.0.1:80/v1/querySmartContract' \
-H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"contractIndex": 11,
"inputType": 12,
"inputSize": 0,
"requestData": ""
}'

## 13. getInvocator (no input, id(32 bytes) output)

curl -X 'POST' \
'http://127.0.0.1:80/v1/querySmartContract' \
-H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"contractIndex": 11,
"inputType": 13,
"inputSize": 0,
"requestData": ""
}'

## 14. getTotalLockedTokens (no input, uint64 ouput)

curl -X 'POST' \
'http://127.0.0.1:80/v1/querySmartContract' \
-H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"contractIndex": 11,
"inputType": 14,
"inputSize": 0,
"requestData": ""
}'
---

### Private internal methods
## 6. isAdmin
## 7. isManager