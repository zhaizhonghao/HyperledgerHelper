# HyperledgerHelper
It is a tool to help users automatically generate neccesary materials for building a hyperledger fabric network. 
The materails include:
* crypto-config.yaml
* configtx.yaml
* `crypto-config` folder which contains identities  materials of participants in the network
* `genesis.block` of the system channel (the orderers channel)
* `[channelName].tx` which is the transaction for creating a specific channel
* docker-compose.yaml for fabric CA, orderers and peers.

It is also a tool to make users up a network and deploy the smart contract easier. In detail, it can help users 
* create a channel with `[channelName].tx`
* join all the peers into the channel 
* install a chaincode
* approve a chaincode 
* check the aprovals of a chaincode
* instantiate a chaincode
* initialize a chaincode
* invoke a chaincode

## Usage
Compile the tool to generate a executable file `configtxTool`
```
go build
```
Run the tool
```
./configToll
```
Once the tool is running, it will listen on port 8181 of localhost. The tool will listen for HTTP request from front-end or Postman.

## RESTful API
### To generate crypto-config.yaml, `crypto-config` folder and docker-compose.yaml
**Method**:POST

**Content-type**:application/json

**Endpoint**:
```
http://localhost:8181/crypto
```

**Request body payload**:

|Property Name|Type|Description|
|----|----|----|
|OrdererCps|[]OrdererCp|the information of orderers|
|PeerOrgCps|[]PeerOrgCp|the information of Organizations|
**Response payload**:
|title|title|
|----|----|
|item|item|
---
### To generate configtx.yaml, genesis.block of sys channel, [channelName].tx
**Method**:

**Content-type**:

**Endpoint**:
```
http://
```

**Request body payload**:

|title|title|
|----|----|
|item|item|
**Response payload**:
|title|title|
|----|----|
|item|item|
---
### To start all nodes of the network up
**Method**:

**Content-type**:

**Endpoint**:
```
http://
```

**Request body payload**:

|title|title|
|----|----|
|item|item|
**Response payload**:
|title|title|
|----|----|
|item|item|
---
### To create a channel
**Method**:

**Content-type**:

**Endpoint**:
```
http://
```

**Request body payload**:

|title|title|
|----|----|
|item|item|
**Response payload**:
|title|title|
|----|----|
|item|item|
---
### To join all the peer into the channel
**Method**:

**Content-type**:

**Endpoint**:
```
http://
```

**Request body payload**:

|title|title|
|----|----|
|item|item|
**Response payload**:
|title|title|
|----|----|
|item|item|
---
### To package the chaincode
**Method**:

**Content-type**:

**Endpoint**:
```
http://
```

**Request body payload**:

|title|title|
|----|----|
|item|item|
**Response payload**:
|title|title|
|----|----|
|item|item|
---
### To install the chaincode
**Method**:

**Content-type**:

**Endpoint**:
```
http://
```

**Request body payload**:

|title|title|
|----|----|
|item|item|
**Response payload**:
|title|title|
|----|----|
|item|item|
---
### To fetch the package id of the chaincode
**Method**:

**Content-type**:

**Endpoint**:
```
http://
```

**Request body payload**:

|title|title|
|----|----|
|item|item|
**Response payload**:
|title|title|
|----|----|
|item|item|
---
### To approve a chaincode
**Method**:

**Content-type**:

**Endpoint**:
```
http://
```

**Request body payload**:

|title|title|
|----|----|
|item|item|
**Response payload**:
|title|title|
|----|----|
|item|item|
---
### To check the approvals of a chaincode
**Method**:

**Content-type**:

**Endpoint**:
```
http://
```

**Request body payload**:

|title|title|
|----|----|
|item|item|
**Response payload**:
|title|title|
|----|----|
|item|item|
---
### To instantiate a chaincode
**Method**:

**Content-type**:

**Endpoint**:
```
http://
```

**Request body payload**:

|title|title|
|----|----|
|item|item|
**Response payload**:
|title|title|
|----|----|
|item|item|
---
### To initialize a chaincode
**Method**:

**Content-type**:

**Endpoint**:
```
http://
```

**Request body payload**:

|title|title|
|----|----|
|item|item|
**Response payload**:
|title|title|
|----|----|
|item|item|

# ?????????????????????????????????
## ??????
????????????????????????????????????????????????????????????????????????????????????
* crypto-config.yaml
* configtx.yaml
* ?????????crpto-config.yaml??????????????????????????????ca???orderer???peer???couchDB???docker-compose.yaml

??????????????????
![2c5c3c759a4017d2756b090754aef4b4.png](en-resource://database/13414:1)


????????????
**crypto-config.yaml**
????????????????????????????????????????????????Decentralized identity management????????????????????????????????????????????????????????????????????????
????????????????????????????????????????????????????????????????????????????????????
???1?????????????????????????????????Fabric-ca?????????????????????
???2????????????????????????????????????????????????????????????*cryptogen*???????????????*crypto-config.yaml*???????????????????????????????????????????????????
crypto-config.yaml??????????????????????????????OrdererOrgs???PeerOrgs?????????
OrdererOrgs?????????????????????orderer?????????
PeerOrgs????????????????????????????????????????????????
