package chain
import (
	"log"
	"context"
	"math/big"
	"github.com/golang/glog"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"encoding/json"
	"crypto/ecdsa"
	"net/http"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"bytes"
	//"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	//"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	newi "github.com/onecool2/web-server/contract" // for demo
	"time"
	"sync"
	"fmt"
	"strings"
)
const (
	RPC_HOST string = "http://115.159.19.208:32004"
	WS_HOST string = "ws://115.159.19.208:32001"
	CONTRACT_ADDRESS string = "0xfA02a776BB22cc644AE4d78EC348702bFB5D927A"
	OWNER_PUBLIC_KEY string = "0xa00dd4406d2dd1d8fde543e2150203ae701e4701"
	OWNER_PRIVATE_KEY string = "41177460dc4d1760832fcbfcaae19e3f734c968edafd00e2947e8bee3da59801"
	EMPTY int = 0
	READY int = 1
	SENT int = 2
	MAX_ON_CHAIN_DELAY = "-30s"
)
/*
const (
	RPC_HOST string = "http://111.230.101.228:32000"
	WS_HOST string = "ws://111.230.101.228:32001"
	CONTRACT_ADDRESS string = "0x01c08d2A8F27702f1aC6A85A7f14064653671636"
	OWNER_PUBLIC_KEY string = "0xef86b1d1eb61f7a817f6b7c21d4363d2bc46fa65"
	OWNER_PRIVATE_KEY string = "1e0f1edd98830544546714e85d18fa4d90cebe7600dec5d6d43886a680c1175b"
	EMPTY int = 0
	READY int = 1
	SENT int = 2
	MAX_ON_CHAIN_DELAY = "-30s"
)*/
/*************************************************************************/

/*************************************************************************/

type LogTransfer struct {
    From   common.Address
    To     common.Address
    Value  *big.Int
    Time   *big.Int
}

type Transfer struct {
	Method   common.Address
	To     common.Address
	Value  *big.Int
}

type SenderBuffer struct {
	Function string
	Arg[4] string
	tx common.Hash
	state int
	time time.Time
	m *sync.Mutex
}

var (
	EthClient *ethclient.Client
	Contract *newi.Newi
	privateKey *ecdsa.PrivateKey
	publicKey common.Address
	owner common.Address
)

var senderQueue [3000]SenderBuffer

func Insert(function string, arg[4] string, tx common.Hash) (int){
	for i := 0; i < len(senderQueue); i++ {
		senderQueue[i].m.Lock()
		if senderQueue[i].state == EMPTY {
			senderQueue[i].Arg[0] = arg[0]
			senderQueue[i].Arg[1] = arg[1]
			senderQueue[i].Arg[2] = arg[2]
			senderQueue[i].Arg[3] = arg[3]
			senderQueue[i].state = SENT
			senderQueue[i].Function = function
			senderQueue[i].tx = tx
			/*p.dataBuffer[i].tx == 0
			fmt.Println()
			fmt.Println("insert:", data)
			fmt.Println("i:", i)
			fmt.Println("p.dataBuffer[%d].state", i, p.dataBuffer[i].state)
			*/
			fmt.Println("insert senderQueue:", i)
			senderQueue[i].m.Unlock()
			return 0
		}
		senderQueue[i].m.Unlock()
	}
	return -1
}

func LoopAndSendTx(){
	for {
		select {
			case	senderBuffer := <-SendChan:
			if senderBuffer.Function == "transfer" {
				senderBuffer.tx = transferToken(senderBuffer.Arg[0], senderBuffer.Arg[1], senderBuffer.Arg[2])
			}else if senderBuffer.Function == "allocatedToken" {
				senderBuffer.tx = allocateTokens(senderBuffer.Arg[0], senderBuffer.Arg[1])
			}else {
				log.Fatal("found unsupport function:", senderBuffer.Function)
			}
			err := Insert(senderBuffer.Function, senderBuffer.Arg, senderBuffer.tx)
			if err != 0 {
				glog.Warning ("Insert failed, senderQueue is full")
			}
			//fmt.Println("insert:", .Function)
			//fmt.Println("args:", senderQueue[i].Arg[0], senderQueue[i].Arg[1], senderQueue[i].Arg[2], senderQueue[i].Arg[3])
			//fmt.Println("tx:", senderQueue[i].tx)
		}
		fmt.Println("get a tx")
	}
}

func LoopAndRemove(tx common.Hash){
	now := time.Now()
        m, _ := time.ParseDuration(MAX_ON_CHAIN_DELAY)
	startPoint := now.Add(m)
	for i := 0; i < len(senderQueue); i++ {
		senderQueue[i].m.Lock()
		if senderQueue[i].state == SENT {
			if senderQueue[i].tx == tx {
				senderQueue[i].state = EMPTY
				fmt.Println("Remove senderQueue:", i)
				fmt.Println("remove:", senderQueue[i].Function)
				fmt.Println("args:", senderQueue[i].Arg[0], senderQueue[i].Arg[1], senderQueue[i].Arg[2], senderQueue[i].Arg[3])
				fmt.Println("tx:", senderQueue[i].tx)
			}else {
				if senderQueue[i].time.Before(startPoint) {
					glog.Warning("!!!tx haven't found util now", senderQueue[i].tx)
					//should re-sent
				}
			}
		}
		senderQueue[i].m.Unlock()
	}
}
/********************************************************************************************************/
/*********In order to resove "import cycle not allowed", move this function from sidecarproxy.go ********/
/********************************************************************************************************/
type Event struct {
    From string `json: "from"`
    To   string `json: "to"`
    Value string `json: "value"`
}


func SendEventToZrServer(from string, to string, value string) {
    var event Event
    event.Value = value
    event.From = from
    event.To = to
    fmt.Println("from to value", from, to, value)
    if bs, err := json.Marshal(event); err == nil {
	fmt.Println("send to ZrServer:", string(bs))
        req := bytes.NewBuffer([]byte(bs))
	fmt.Println("send to ZrServer:", req)
        //tmp := `{"name":"junneyang", "age": 88}`
        //req = bytes.NewBuffer([]byte(tmp))

        body_type := "Content-Type: application/json"
	/*resp, _ := */http.Post("http://127.0.0.1:3001/event", body_type, req)
        //body, _ := ioutil.ReadAll(resp.Body)
        //fmt.Println(string(body))
    } else {
        fmt.Println(err)
    }
}
/********************************************************************************************************/

func GoThroughBlock(){
	var currentBlockNum *big.Int
	chainID, err := EthClient.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	currentBlock, err := EthClient.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	currentBlockNum = currentBlock.Number()

	hash := sha3.NewKeccak256()
	transferFnSignature := []byte("transfer(address,uint256)")
	hash.Write(transferFnSignature)
	transferMethod := hexutil.Encode(hash.Sum(nil)[:4])
	for {
		latestBlock, err := EthClient.BlockByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
		if currentBlockNum.Uint64() < latestBlock.Number().Uint64() {
			for currentBlockNum.Uint64() < latestBlock.Number().Uint64() {
				currentBlockNum.Add(currentBlockNum, big.NewInt(1))
				block, err := EthClient.BlockByNumber(context.Background(), currentBlockNum)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("/#####################################/")        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
				fmt.Println(block.Number().Uint64())     // 5671744
				//fmt.Println(block.Time().Uint64())       // 1527211625
				fmt.Println(len(block.Transactions()))   // 144
				for _, tx := range block.Transactions() {
					glog.Info("/********************************/")        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
					glog.Info("tx hash:", tx.Hash().Hex())        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
					glog.Info("tx value:", tx.Value().String())    // 10000000000000000
					glog.Info("tx gas:", tx.Gas())               // 105000
					glog.Info("tx gasPrice:", tx.GasPrice().Uint64()) // 102000000000
					glog.Info("tx nonce:", tx.Nonce())             // 110644
					glog.Info("tx data:", tx.Data()[16:35])              // []
					if tx.To() != nil {
						glog.Info("tx to:", tx.To().Hex())          // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e					
					}
					if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err == nil {
						glog.Info("msg from:", msg.From().Hex()) // 0x0fD081e3Bb178dc45c0cb23202069ddA57064258
					}
					receipt, err := EthClient.TransactionReceipt(context.Background(), tx.Hash())
					if err != nil {
						log.Fatal(err)
					}

					glog.Info("receipt status:", receipt.Status)
					glog.Info("/********************************/")        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
					if receipt.Status == 1 {
						if tx.To().Hex() == CONTRACT_ADDRESS {
							method := hexutil.Encode(tx.Data()[:4])
							if method == transferMethod {
								if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err == nil {
									LoopAndRemove(tx.Hash())
									to := hexutil.Encode(tx.Data()[16:36])
									value := hexutil.Encode(tx.Data()[36:68])
									fmt.Println("value:", value)              // []
									val, err1 := new(big.Int).SetString(value[2:], 16)
									if err1 != true {
										log.Fatal(err1)
									}              // []
									fmt.Println("val:", val)              // []
									glog.Info("method:", hexutil.Encode(tx.Data()[:4]))              // []
									fmt.Println("to:", hexutil.Encode(tx.Data()[16:36]))              // []
									glog.Info("from:", msg.From().Hex())
									glog.Info("OWNER_PUBLIC_KEY:",  OWNER_PUBLIC_KEY)
									fmt.Printf("from:%b\n", msg.From().Hex())
									fmt.Printf("OWNER_PUBLIC_KEY:%b\n", OWNER_PUBLIC_KEY)
									ff:=msg.From().Hex()
									fmt.Println(strings.Compare(ff, OWNER_PUBLIC_KEY))
									if 0 != strings.Compare(msg.From().Hex(), OWNER_PUBLIC_KEY) {
										SendEventToZrServer(msg.From().Hex(), to, val.String())
									}
								}
							}
						}
					}
				}
				fmt.Println("/#####################################/")        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
			}
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func allocateTokens(to string, value string) (common.Hash) {
    var toArray []common.Address
    var valueArray []*big.Int
    fmt.Println("allocateTokens:", to)
    _ = Contract

    nonce, err := EthClient.PendingNonceAt(context.Background(), owner)
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

    tx, err := Contract.AllocateTokens(auth, toArray, valueArray)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("send token:", to, value)
    fmt.Println("tx sent: %s", tx.Hash().Hex()) // tx sent: 0x8d490e535678e9a24360e955d75b27ad307bdfb97a1dca51d0f3035dcee3e870
    return tx.Hash()
}

func transferToken(from string, to string, value string)(common.Hash) {

    fmt.Println("transferToken:", from , to)
    _ = Contract

    nonce, err := EthClient.PendingNonceAt(context.Background(), owner)
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
    //wxx_from := common.HexToAddress(from)

    //_value := big.NewInt(10)
    _value, _ := new(big.Int).SetString(value, 10)
    //if err1 != nil {
    //    log.Fatal(err1)
    //}
    //copy(_to[:], address("0x"))
    //copy(_value[:], uint256(10))

 /*  wxx tx, err := Contract.TransferFrom(auth, _from, _to, _value)
    if err != nil {
	    glog.Warning("send tx:", err)
    }*/
    tx, err := Contract.Transfer(auth, _to, _value)
    if err != nil {
	    glog.Warning("send tx:", err)
    }
    fmt.Printf("send token:", from, to, value)
    fmt.Printf("tx sent: %s", tx.Hash().Hex()) // tx sent: 0x8d490e535678e9a24360e955d75b27ad307bdfb97a1dca51d0f3035dcee3e870
    fmt.Println("")
    return tx.Hash()
    /*result, err := contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress("0xfA02a776BB22cc644AE4d78EC348702bFB5D927A"))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("result", result) // "1.0"
    */
}

var SendChan chan SenderBuffer
func init(){
	var contractAddress common.Address
	var err error
	EthClient, err = ethclient.Dial(RPC_HOST)
	if err != nil {
		log.Fatal(err)
	}

	contractAddress = common.HexToAddress(CONTRACT_ADDRESS)
	Contract, err = newi.NewNewi(contractAddress, EthClient)
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
	/******************************************************/
	fmt.Printf("init-senderQueue")
	for i := 0; i < len(senderQueue); i++ {
		senderQueue[i].state = EMPTY
		senderQueue[i].m = new(sync.Mutex)
	}


        SendChan = make(chan SenderBuffer, 100)
	/*
	 //event
    wsclient, err := Ethclient.Dial(WS_HOST)
    if err != nil {
        log.Fatal(err)
    }

    address := common.HexToAddress(CONTRACT_ADDRESS)
    query := ethereum.FilterQuery{
        Addresses: []common.Address{address},
    }
    //contract, err := newi.NewNewi(address, wsclient)
    //if err != nil {
    //    log.Fatal(err)
    //}
    logs := make(chan types.Log)
    sub, err := wsclient.SubscribeFilterLogs(context.Background(), query, logs)
    if err != nil {
            log.Fatal(err)
    }
    watchBlockChain()
    fmt.Println("Web service is listening...")
    for {
        select {
        case err := <-sub.Err():
            log.Fatal(err)
        case vLog := <-logs:
            decode_event(contract, vLog) // pointer to event log
        }
    }
*/

}

/*
func decode_event(Contract *newi.Newi, logs types.Log) {
    fmt.Println(logs) // pointer to event log
    fmt.Println("topic0:%x", logs.Topics[0].Hex())
    fmt.Println("topic1:%x", logs.Topics[1].Hex())
    fmt.Println("topic2:%x", logs.Topics[2].Hex())
    contractAbi, err := abi.JSON(strings.NewReader(string(newi.NewiABI)))
    if err != nil {
        log.Fatal(err)
    }
    //logTransferSig := []byte("Transfer(address,address,uint256)")
    logTransferSig := []byte("Transfer(address, address, uint256)")
    logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

    //for _, vLog := range logs {
        fmt.Printf("Log Block Number: %d\n", logs.BlockNumber)
        fmt.Printf("Log Block Number: %d\n", logs.BlockNumber)
        fmt.Printf("Log Index: %d\n", logs.Index)
        fmt.Printf("Log Topics: %x\n", logs.Topics[0])
        fmt.Printf("Expect: %x\n", logTransferSigHash)

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

            fmt.Println("from:", transferEvent.From.Hex())
            fmt.Println("to:", transferEvent.To.Hex())
            fmt.Println("value:", transferEvent.Value.String())
	    fmt.Println("time:", transferEvent.Time.String())
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

	    fmt.Println("from:", transferEvent.From.Hex())
            fmt.Println("to:", transferEvent.To.Hex())
            fmt.Println("value:", transferEvent.Value.String())
	    fmt.Println("time:", transferEvent.Time.String())

	    from := transferEvent.From.String()
	    to := transferEvent.To.String()
	    value := transferEvent.Value.String()
	    if (from != OWNER_PUBLIC_KEY) {
		sendEventToZrServer(from, to, value)
            }
        }

        fmt.Printf("\n\n")
}
*/
