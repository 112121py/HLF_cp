
## Chaincode
| 模組名稱 | Chaincode | 功能描述 |
|----------|-----------|----------|
| 任務管理合約 | `FLTaskContract` | 創建、管理 FL 任務（任務類型、輪次、參與者、參數） |
| 模型狀態合約 | `ModelContract` | 模型提交、驗證、貢獻紀錄、版本追蹤 |
| 全局資訊合約 | `ChannelStatsContract` | 查詢 channel 狀態（資料量、模型次數、組織貢獻等） |
| 模型背書合約 | `ZeroTrustEndorseContract` | 驗證 hash、簽名、來源可信度、模型完整性 |
| （可選）模型分享合約 | `ModelExchangeContract` | 模型下載連結管理與快取驗證 |

## 開發與部署步驟（Go）
- `<GO_MODULE_NAME>.go` 合約定義
- `main.go`（必要，負責註冊合約）
- `go.mod`（初始化即可），進行 Go module 定義
    ```
    go mod init <GO_MODULE_NAME>
    go get github.com/hyperledger/fabric-contract-api-go/contractapi
    ```  
- `go.sum`
  - 用 `go mod tidy` 幫你補全




## 生成 channel block
```
export FABRIC_CFG_PATH=$PWD/configtx
mkdir -p channel-artifacts

configtxgen \
  -profile TwoOrgsApplicationGenesis \
  -channelID flchannel \
  -outputBlock ./channel-artifacts/flchannel.block
```


## 節點加入 channel
### orderer 節點加入
```
osnadmin channel join \
  --channelID flchannel \
  --config-block ./channel-artifacts/flchannel.block \
  --ca-file ./organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt \
  --client-cert ./organizations/ordererOrganizations/example.com/users/Admin@example.com/msp/signcerts/cert.pem \
  --client-key ./organizations/ordererOrganizations/example.com/users/Admin@example.com/msp/keystore/<your_sk>.pem \
  --orderer-address localhost:7053
```
如果成功，會出現以下結果
```
Status: 201
{
	"name": "flchannel",
	"url": "/participation/v1/channels/flchannel",
	"consensusRelation": "consenter",
	"status": "active",
	"height": 1
}
```
### peer 節點加入
在本機將 channel block 傳入 org 容器，成功的話會看到 `Successfully copied 20.5kB to peer0.org3.example.com:/tmp/flchannel.block` 之類的訊息
```
docker cp ./channel-artifacts/flchannel.block peer0.org3.example.com:/tmp/flchannel.block
```
#### ~~不要以 org3 為例子~~

#### 以 org1 為例子
```
docker exec -it peer0.org1.example.com bash
```
```
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
export CORE_PEER_MSPCONFIGPATH=etc/hyperledger/fabric/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
export FABRIC_CFG_PATH=/etc/hyperledger/fabric
```
> 要注意CORE_PEER_MSPCONFIGPATH不要設成/etc/hyperledger/fabric/org1.example.com/msp/
因為這樣你使用的身分就是peer，權限不夠，要換成admin的身分，如果找不到各種身分的msp路徑可以去fabric-samples/zta-fl-network/test-network/organizations/fabric-ca/Enrolled.sh看

> error message: 
![image](https://github.com/user-attachments/assets/59687b6e-1729-452f-b7de-e31d577c1f14)

加入 channel
```
peer channel join -b /tmp/flchannel.block
```
![alt text](../../../src/img/image-3.png)

> 為什麼 peer0.org3.example.com 可當作自己容器內的位址？
> 在 Docker 網路中，每個容器名稱會被登錄進內部 DNS。所以即使你在 peer0.org3 的容器內，peer0.org3.example.com 會解析為本機 loopback（類似你叫自己的 DNS 名）。
> 這也避開了 127.0.0.1 對應不到 cert 的問題（x509 TLS 驗證）。

### 確認是否加入成功?
進入到該組織的 peer 節點容器中執行以下指令
```
peer channel getinfo -c flchannel
```
#### org1
![alt text](../../../src/img/image-4.png)




#### org3
> 註記: 由於 configtx 設定的 TwoOrgsApplicationGenesis 預設只有 org1 和 org2，也就代表了這個通道的 channel config 沒有加入除了 org1 和 org2 以外的組織，因此當一個 channel 尚未將某組織（Org3）加入 channel config 時，即使 peer 成功加入 channel block，也不能查詢資訊、無法 invoke/query chaincode。
```
$ peer channel getinfo -c flchannel

2025-04-07 09:00:02.695 UTC 0001 INFO [channelCmd] InitCmdFactory -> Endorser and orderer connections initialized
Error: received bad response, status 500: access denied for [GetChainInfo][flchannel]: [Failed evaluating policy on signed data during check policy on channel [flchannel] with policy [/Channel/Application/Readers]: [implicit policy evaluation failed - 0 sub-policies were satisfied, but this policy requires 1 of the 'Readers' sub-policies to be satisfied]]
```

## 查看是否 peer 加入任何 channel
### 使用 peer 節點查看
進入 CLI 容器：
```
docker exec -it cli bash
```
查詢 Org1 的 Peer channel
```
CORE_PEER_ADDRESS=peer0.org1.example.com:7051 \
CORE_PEER_LOCALMSPID=Org1MSP \
CORE_PEER_TLS_ENABLED=true \
CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt \
CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/msp \

peer channel list
```
結果範例
```
peer channel list
2025-04-07 08:50:50.197 UTC 0001 INFO [channelCmd] InitCmdFactory -> Endorser and orderer connections initialized
Channels peers has joined: 
flchannel
```
![alt text](../../../src/img/image.png)

### 使用 `osnadmin channel list` 查看 orderer 節點加入哪些 channel
```
osnadmin channel list \
  --ca-file /path/to/orderer/tls/ca.crt \
  --client-cert /path/to/orderer/admin/tls/signcerts/cert.pem \
  --client-key /path/to/orderer/admin/tls/keystore/key.pem \
  --orderer-address localhost:7053
```
範例
```
osnadmin channel list \
  --orderer-address localhost:7053 \
  --ca-file ./organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt \
  --client-cert ./organizations/ordererOrganizations/example.com/users/Admin@example.com/msp/signcerts/cert.pem \
  --client-key ./organizations/ordererOrganizations/example.com/users/Admin@example.com/msp/keystore/<sk.pem>

```
結果範例
```
Status: 200
{
	"systemChannel": null,
	"channels": null
}
```
```
Status: 200
{
	"systemChannel": null,
	"channels": [
		{
			"name": "flchannel",
			"url": "/participation/v1/channels/flchannel"
		}
	]
}
```
![alt text](../../../src/img/image-2.png)

## 在 channel 上使用 chaincode
```
docker cp fltask.tar.gz peer0.org1.example.com:/opt/gopath/src/github.com/hyperledger/fabric/peer/
```


## chaincode
### 打包 chaincode
```
peer lifecycle chaincode package <PACKAGE_NAME>.tar.gz \
  --path <FULL_PATH_TO_YOUR_CHAINCODE_FOLDER> \
  --lang golang \
  --label <PACKAGE_LABEL>

```
例如
```
peer lifecycle chaincode package fltask.tar.gz   --path /root/go/src/github.com/ki225/fabric-samples/zta-fl-network/chaincode/FLTaskContract/   --lang golang   --label fltask_1
```
### 安裝 chaincode
在本機把 fltask.tar.gz 丟到容器
```
docker cp fltask.tar.gz peer0.org3.example.com:/opt/gopath/src/github.com/hyperledger/fabric/peer/
```
```
Successfully copied 10.2kB to peer0.org3.example.com:/opt/gopath/src/github.com/hyperledger/fabric/peer/
```
進入到容器安裝 chaincode
```
peer lifecycle chaincode install fltask.tar.gz
```
```
2025-04-07 12:43:14.739 UTC 0043 INFO [lifecycle] InstallChaincode -> Successfully installed chaincode with package ID 'fltask_1:bfd178828f3fa261b38f957fe00be5ef3ee624566709648515fb3c9bf7b7a166'
```
![alt text](../../../src/img/image-5.png)

### Approve chaincode for Org1
跟 Orderer 通訊，把 chaincode 定義送到 Orderer 進行 共識交易（Ordering），如果目前peer0.org1容器還沒有 orderer 的資訊，可以從本機複製過去
```
docker cp ./organizations/ordererOrganizations/example.com/orderers/   peer0.org1.example.com:/etc/hyperledger/fabric/
```
> 用 --tls --cafile 就是在告訴 peer：「我想透過 TLS 與 orderer.example.com:7050 建立安全連線，請用這個 CA 來驗證對方的憑證」
```
peer lifecycle chaincode approveformyorg \
  --channelID flchannel \
  --name fltask \
  --version 1.0 \
  --sequence 1 \
  --package-id <PACKAGE_ID> \
  --tls \
  --cafile /etc/hyperledger/fabric/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  --orderer orderer.example.com:7050
```
![alt text](../../../src/img/image-6.png)
從 docker logs 也可以看到 ![alt text](../../../src/img/image-7.png)


### commit
在 approve 過後需要 commit
```
docker cp \
  ./organizations/peerOrganizations/org1.example.com/peers/ \
  peer0.org2.example.com:/etc/hyperledger/fabric/msp/users/
```
因為在 `configtx.yaml` 裡面設定了 channel 的預設政策，包括 endorsement policy (/Channel/Application/Endorsement)，所以需要用 `--tlsRootCertFiles` 以加入 Org1 和 Org2 的 peer node 都參與 commit 指令，否則會有 `Error: transaction invalidated with status (ENDORSEMENT_POLICY_FAILURE)` 錯誤

```
peer lifecycle chaincode commit \
  --channelID flchannel \
  --name fltask \
  --version 1.0 \
  --sequence 1 \
  --tls \
  --cafile /etc/hyperledger/fabric/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  --orderer orderer.example.com:7050 \
  --peerAddresses peer0.org1.example.com:7051 \
  --tlsRootCertFiles /etc/hyperledger/fabric/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt  \
  --peerAddresses peer0.org2.example.com:9051 \
  --tlsRootCertFiles /etc/hyperledger/fabric/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt 
```
```
2025-04-07 14:17:02.102 UTC 0001 INFO [chaincodeCmd] ClientWait -> txid [5eac39d271ad26bdbdad4e26adf2a843ba183c7cbe04ec28a034b7c80cefe86e] committed with status (VALID) at peer0.org1.example.com:7051
2025-04-07 14:17:02.106 UTC 0002 INFO [chaincodeCmd] ClientWait -> txid [5eac39d271ad26bdbdad4e26adf2a843ba183c7cbe04ec28a034b7c80cefe86e] committed with status (VALID) at peer0.org2.example.com:9051
```

![alt text](../../../src/img/image-8.png)

> 當需要與 orderer 溝通（如 approve、commit、channel create/join）時，一定要用 Orderer TLS CA 憑證當作 --cafile！
### 測試 chaincode
```
peer chaincode invoke -C flchannel -n fltask \
  -c '{"function":"FLTaskContract:CreateTask","Args":["task1", "第一個任務", "org1", "FedAvg"]}'
```
![alt text](../../../src/img/image-9.png)
從 docker log看到

![alt text](../../../src/img/image-10.png)

## 查看 channel 上 chaincode

> 沒有設置 `CORE_PEER_MSPCONFIGPATH` 會出現以下錯誤
> Error: query failed with status: 500 - Failed to authorize invocation due to failed ACL check: Failed verifying that proposal's creator satisfies local MSP principal during channelless check policy with policy [Admins]: [The identity is not an admin under this MSP [Org1MSP]: The identity does not contain OU [ADMIN], MSP: [Org1MSP]]
>

```
export CORE_PEER_MSPID=Org1MSP
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
export CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/org1.example.com/users/Admin@org1.example.com/msp


peer lifecycle chaincode queryinstalled
```
範例輸出
```
$ peer lifecycle chaincode queryinstalled
Installed chaincodes on peer:
Package ID: fltask_1:bfd178828f3fa261b38f957fe00be5ef3ee624566709648515fb3c9bf7b7a166, Label: fltask_1
```


```
peer lifecycle chaincode querycommitted --channelID <channel_name> --name <chaincode_name> \
  --tls --cafile <orderer_tls_ca> --peerAddresses <peer_address> --tlsRootCertFiles <peer_tls_cert>
```
```
peer lifecycle chaincode querycommitted --channelID flchannel --name fltask \
  --tls --cafile /etc/hyperledger/fabric/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  --orderer orderer.example.com:7050 \
  --peerAddresses peer0.org1.example.com:7051 \
  --tlsRootCertFiles /etc/hyperledger/fabric/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt 
```

```
peer lifecycle chaincode approveformyorg \
  --channelID flchannel \
  --name fltask \
  --version 1.0 \
  --sequence 1 \
  --package-id fltask_1:bfd178828f3fa261b38f957fe00be5ef3ee624566709648515fb3c9bf7b7a166 \
  --orderer orderer.example.com:7050 \
  --tls \
  --cafile /etc/hyperledger/fabric/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  --peerAddresses peer0.org1.example.com:7051 \
  --tlsRootCertFiles /etc/hyperledger/fabric/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

```

```
peer channel fetch config flchannel.block \
  -o orderer.example.com:7050 \
  --tls \
  --cafile /etc/hyperledger/fabric/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -c flchannel
```

```
peer channel fetch config ./channel-artifacts/flchannel.block \
  -o orderer.example.com:7050 \
  --tls \
  --cafile $PWD/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -c flchannel
```

