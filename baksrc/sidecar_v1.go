package main

import (
    "log"
    "context"
    "math/big"
    "encoding/json"
    "crypto/ecdsa"
    "strings"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/core/types"
    newi "github.com/onecool2/web-server/contract" // for demo
    "flag"
    "github.com/golang/glog"
    //"matrix-backend/pkg/api"
    "bytes"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "io/ioutil"
    //"html"
)

/************************************************************************
var ZUI_RI_SERVER_HOST string = "http://127.0.0.1:3001/event"
var RPC_HOST string = "http://115.159.19.208:32004"
var WS_HOST string = "ws://115.159.19.208:32003"
var CONTRACT_ADDRESS string = "0xfA02a776BB22cc644AE4d78EC348702bFB5D927A"
var OWNER_PUBLIC_KEY string = "0xa00dd4406d2dd1d8fde543e2150203ae701e4701"
var OWNER_PRIVATE_KEY string = "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
************************************************************************/

/*************************************************************************/
var ZUI_RI_SERVER_HOST string = "http://127.0.0.1:3001/event"
var RPC_HOST string = "http://111.230.101.228:32000"
var WS_HOST string = "ws://111.230.101.228:32001"
var CONTRACT_ADDRESS string = "0xd492825cA3427d4D39744C010096803132d54B53"
var OWNER_PUBLIC_KEY string = "0xef86b1d1eb61f7a817f6b7c21d4363d2bc46fa65"
var OWNER_PRIVATE_KEY string = "1e0f1edd98830544546714e85d18fa4d90cebe7600dec5d6d43886a680c1175b"
/*************************************************************************/


/*************************************************************************/
var ethClient *ethclient.Client
var contract *newi.Newi
var privateKey *ecdsa.PrivateKey
var publicKey common.Address
var owner common.Address
/*************************************************************************/
type RetMsg struct {
    code int //400
    message string //"具体错误信息",
    data string //"data":null
}

func replyMsg (writer http.ResponseWriter, code int, message string, data string){

    var retMsg RetMsg
    retMsg.code = code
    retMsg.message = message
    retMsg.data = data
    if bs, err :=json.Marshal(retMsg); err == nil {
        //req := bytes.NewBuffer([]byte(bs))
        //body_type := "application/x-www-form-urlencoded"
	//writer.Header().Set("Content-type", "application/text")
	writer.Header().Set("Content-type", "application/x-www-form-urlencoded")
	writer.WriteHeader(200)
        writer.Write([]byte(data))

        //fmt.Fprintln(writer, bs)
        /*resp, _ := *///http.Post(ZUI_RI_SERVER_HOST, body_type, req)
        fmt.Println(bs)
    }

}

func transferHandler(writer http.ResponseWriter, request *http.Request) {
    body, _:= ioutil.ReadAll(request.Body)
    defer request.Body.Close()
    body_str := string(body)
    fmt.Println(body_str)
    var transfer map[string]interface{}
    if err := json.Unmarshal(body, &transfer); err == nil {
        fmt.Println(transfer["serialNumber"])
        fmt.Println(transfer["from"])
        fmt.Println(transfer["to"])
        fmt.Println(transfer["value"])

	serialNumber := transfer["serialNumber"].(string)
	from := transfer["from"].(string)
        to := transfer["to"].(string)
        value := transfer["value"].(string)

	transferToken(from, to , value, serialNumber)  //fmt.Fprint(w, string(ret))
    } else {
	fmt.Println("Unmarshal:", err)
    }
}

func allocateHandler(writer http.ResponseWriter, request *http.Request) {
    body, _:= ioutil.ReadAll(request.Body)
    defer request.Body.Close()
    body_str := string(body)
    fmt.Println(body_str)
    var transfer map[string]interface{}
    if err := json.Unmarshal(body, &transfer); err == nil {
        fmt.Println(transfer["to"])
        fmt.Println(transfer["value"])

        to := transfer["to"].(string)
        value := transfer["value"].(string)

	allocateTokens(to , value)  //fmt.Fprint(w, string(ret))
    } else {
	fmt.Println("Unmarshal:", err)
    }
}
func balanceHandler(writer http.ResponseWriter, request *http.Request) {
    body, _:= ioutil.ReadAll(request.Body)
    defer request.Body.Close()
    body_str := string(body)
    fmt.Println(body_str)
    var balance map[string]interface{}
    if err := json.Unmarshal(body, &balance); err == nil {
        fmt.Println(balance["address"])
	address := balance["address"].(string)
	bal := balanceOf(address) //fmt.Fprint(w, string(ret))
        //fmt.Fprintln(writer, "Hello, ", html.EscapeString(ret.String()))
	replyMsg(writer, 200, "success", bal.String())
    } else {
	fmt.Println("Unmarshal:", err)
    }
}

func startServer() {
    r := mux.NewRouter()
    r.HandleFunc("/allocate/tokens", allocateHandler)
    r.HandleFunc("/transfer", transferHandler)
    r.HandleFunc("/balance", balanceHandler)
    http.Handle("/transfer", r)
    http.Handle("/balance", r)
    http.Handle("/allocate/tokens", r)
    fmt.Println("Side Car proxy start listening on 3000...")
    http.ListenAndServe(":3000", r)
}

type Event struct {
    from string
    to   string
    value string
}

func sendEventToZrServer(from string, to string, value string) {
    var event Event
    event.value = value
    event.from = from
    event.to = to
    if bs, err := json.Marshal(event); err == nil {
        fmt.Println("send to ZrServer:", string(bs))
        req := bytes.NewBuffer([]byte(bs))
        //tmp := `{"name":"junneyang", "age": 88}`
        //req = bytes.NewBuffer([]byte(tmp))

        body_type := "Content-Type: application/json"
        /*resp, _ := */http.Post(ZUI_RI_SERVER_HOST, body_type, req)
        //body, _ := ioutil.ReadAll(resp.Body)
        //fmt.Println(string(body))
    } else {
        fmt.Println(err)
    }
}

/*************************************************************************/
type LogTransfer struct {
    From   common.Address
    To     common.Address
    Value  *big.Int
    Time   *big.Int
}


func decode_event(Contract *newi.Newi, logs types.Log) {
    fmt.Println(logs) // pointer to event log
    contractAbi, err := abi.JSON(strings.NewReader(string(newi.NewiABI)))
    if err != nil {
        log.Fatal(err)
    }
    logTransferSig := []byte("Transfer(address,address,uint256)")
    logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

    //for _, vLog := range logs {
        fmt.Printf("Log Block Number: %d\n", logs.BlockNumber)
        fmt.Printf("Log Block Number: %d\n", logs.BlockNumber)
        fmt.Printf("Log Index: %d\n", logs.Index)
        fmt.Printf("Log Topics: %x\n", logs.Topics[0])
        fmt.Printf("Expect: %x\n", logTransferSigHash.Hex())

        switch logs.Topics[0].Hex() {
	case logTransferSigHash.Hex():
            fmt.Printf("Log Name: Transfer\n")

            var transferEvent LogTransfer

            err := contractAbi.Unpack(&transferEvent, "Transfer", logs.Data)
            if err != nil {
                log.Fatal(err)
            }

            transferEvent.From = common.HexToAddress(logs.Topics[1].Hex())
            transferEvent.To = common.HexToAddress(logs.Topics[2].Hex())
            //transferEvent.Value = logs.Topics[3]
            //transferEvent.Time = logs.Topics[4](big.Int)

            fmt.Printf("From: %s\n", transferEvent.From.Hex())
            fmt.Printf("To: %s\n", transferEvent.To.Hex())
            fmt.Printf("Value: %s\n", transferEvent.Value.String())
           // fmt.Printf("Time: %s\n", transferEvent.Time.String())

	    from := transferEvent.From.String()
	    to := transferEvent.To.String()
	    value := transferEvent.Value.String()
	    if (from != "0x96216849c49358B10257cb55b28eA603c874b05E") {
		    sendEventToZrServer(from, to, value)
            }
        default:
	    fmt.Printf("Log Name: Transfer\n")
            var transferEvent LogTransfer

            err := contractAbi.Unpack(&transferEvent, "Transfer", logs.Data)
            if err != nil {
                log.Fatal(err)
            }

            transferEvent.From = common.HexToAddress(logs.Topics[1].Hex())
            transferEvent.To = common.HexToAddress(logs.Topics[2].Hex())

            fmt.Printf("From: %s\n", transferEvent.From.Hex())
            fmt.Printf("To: %s\n", transferEvent.To.Hex())
            fmt.Printf("Value: %s\n", transferEvent.Value.String())
            fmt.Printf("TX: %s\n", logs.TxHash.Hex())

	    from := transferEvent.From.String()
	    to := transferEvent.To.String()
	    value := transferEvent.Value.String()
	    if (from != OWNER_PUBLIC_KEY) {
		sendEventToZrServer(from, to, value)
            }
        }

        fmt.Printf("\n\n")
}

func init() {
    var contractAddress common.Address
    var err error
    ethClient, err = ethclient.Dial(RPC_HOST)
    if err != nil {
        log.Fatal(err)
    }

    contractAddress = common.HexToAddress(CONTRACT_ADDRESS)
    contract, err = newi.NewNewi(contractAddress, ethClient)
    if err != nil {
        log.Fatal(err)
    }

    privateKey, err = crypto.HexToECDSA(OWNER_PRIVATE_KEY)
    if err != nil {
        log.Fatal(err)
    }
    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
    }

    owner = crypto.PubkeyToAddress(*publicKeyECDSA)
}

func balanceOf(address string) *big.Int {

    fmt.Println("call balance")
    balance, err := contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress(address))
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(balance) // "1.0"

    return  balance
}

func allocateTokens(to string, value string) {
    var toArray []common.Address
    var valueArray []*big.Int
    fmt.Println("transferToken:", to)
    _ = contract

    nonce, err := ethClient.PendingNonceAt(context.Background(), owner)
    if err != nil {
        log.Fatal(err)
    }

    auth := bind.NewKeyedTransactor(privateKey)
    auth.Nonce = big.NewInt(int64(nonce))
    auth.Value = big.NewInt(0)     // in wei
    auth.GasLimit = uint64(300000) // in units
    /* remove following lines for quorum blockchain
    gasPrice, err := ethClient.SuggestGasPrice(context.Background())
    auth.GasPrice = gasPrice
    if err != nil {
        log.Fatal(err)
    }*/
    toList := strings.Split(to,",")
    for _, item := range toList {
        toArray = append(toArray, common.HexToAddress(item))
    }

    valueList := strings.Split(value,",")
    for _, item := range valueList {
        value, _ := new(big.Int).SetString(item, 10)
        valueArray = append(valueArray, value)
    }

    tx, err := contract.AllocateTokens(auth, toArray, valueArray)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("send token:", to, value)
    fmt.Println("tx sent: %s", tx.Hash().Hex()) // tx sent: 0x8d490e535678e9a24360e955d75b27ad307bdfb97a1dca51d0f3035dcee3e870
    fmt.Println("")
    /*result, err := contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress("0xfA02a776BB22cc644AE4d78EC348702bFB5D927A"))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("result", result) // "1.0"
    */
}
func transferToken(from string, to string, value string, serialNumber string) {

    fmt.Println("transferToken:", from , to)
    _ = contract

    nonce, err := ethClient.PendingNonceAt(context.Background(), owner)
    if err != nil {
        log.Fatal(err)
    }

    auth := bind.NewKeyedTransactor(privateKey)
    auth.Nonce = big.NewInt(int64(nonce))
    auth.Value = big.NewInt(0)     // in wei
    auth.GasLimit = uint64(300000) // in units
    /* remove following lines for quorum blockchain
    gasPrice, err := ethClient.SuggestGasPrice(context.Background())
    auth.GasPrice = gasPrice
    if err != nil {
        log.Fatal(err)
    }*/

    _to := common.HexToAddress(to)
    _from := common.HexToAddress(from)

    //_value := big.NewInt(10)
    _value, _ := new(big.Int).SetString(value, 10)
    //if err1 != nil {
    //    log.Fatal(err1)
    //}
    //copy(_to[:], address("0x"))
    //copy(_value[:], uint256(10))

    tx, err := contract.TransferFrom(auth, _from, _to, _value)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("send token:", from, to, value)
    fmt.Printf("tx sent: %s", tx.Hash().Hex()) // tx sent: 0x8d490e535678e9a24360e955d75b27ad307bdfb97a1dca51d0f3035dcee3e870
    fmt.Println("")
    /*result, err := contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress("0xfA02a776BB22cc644AE4d78EC348702bFB5D927A"))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("result", result) // "1.0"
    */
}

func main() {
    //init()
    flag.Parse()
    defer glog.Flush()

    go startServer()
    //event
    wsclient, err := ethclient.Dial(WS_HOST)
    if err != nil {
        log.Fatal(err)
    }

    address := common.HexToAddress(CONTRACT_ADDRESS)
    query := ethereum.FilterQuery{
        Addresses: []common.Address{address},
    }
    /*contract, err := newi.NewNewi(address, wsclient)
    if err != nil {
        log.Fatal(err)
    }*/
    logs := make(chan types.Log)
    sub, err := wsclient.SubscribeFilterLogs(context.Background(), query, logs)
    if err != nil {
	    log.Fatal(err)
    }

    fmt.Println("Web service is listening...")
    for {
        select {
        case err := <-sub.Err():
            log.Fatal(err)
        case vLog := <-logs:
            decode_event(contract, vLog) // pointer to event log
        }
    }
}
