package sidecarproxy

import (
    "log"
    //"context"
    "math/big"
    "encoding/json"
    //"crypto/ecdsa"
    //"strings"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    //"matrix-backend/pkg/api"
    "bytes"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "io/ioutil"
    "github.com/ethereum/go-ethereum/common"
    //"time"
    "github.com/onecool2/zuiri/sidecar/pkg/chain"
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

/*************************************************************************/


type SideCarProxy struct {
	HostName string
	Event string
	Balance string
	AllocateTokens string
	Transfer string
}

type RetMsg struct {
	code int //400
	message string //"具体错误信息",
	data string //"data":null
}


func (s *SideCarProxy) replyMsg (writer http.ResponseWriter, code int, message string, data string){

    retMsg := RetMsg{code, message, data}
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

func (s *SideCarProxy) TransferHandler (writer http.ResponseWriter, request *http.Request) {
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

	//serialNumber := transfer["serialNumber"].(string)
	from := transfer["From"].(string)
        to := transfer["To"].(string)
        value := transfer["Value"].(string)
	senderBuffer := chain.SenderBuffer{
		Function: "transfer",
		//		Arg: arg //[4]string{ serialNumber, from, to, value }
	}
	//senderBuffer.Arg[1] = serialNumber
	senderBuffer.Arg[0] = from
	senderBuffer.Arg[1] = to
	senderBuffer.Arg[2] = value

	chain.SendChan <- senderBuffer
    } else {
	s.replyMsg(writer, 400, err.Error(), err.Error())
	fmt.Println("Unmarshal:", err)
    }
    fmt.Println("end TransferHandler")
}

func (s *SideCarProxy) AllocateHandler(writer http.ResponseWriter, request *http.Request) {
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
	senderBuffer := chain.SenderBuffer{
		Function: "allocatedToken",
		//		Arg: arg //[4]string{ serialNumber, from, to, value }
	}
	senderBuffer.Arg[0] = to
	senderBuffer.Arg[1] = value

        chain.SendChan <- senderBuffer

	//s.replyMsg(writer, 200, "success", tx)
    } else {
	s.replyMsg(writer, 400, err.Error(), err.Error())
	fmt.Println("Unmarshal:", err)
    }
}
func (s *SideCarProxy) BalanceHandler(writer http.ResponseWriter, request *http.Request) {
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
	s.replyMsg(writer, 200, "success", bal.String())
    } else {
	s.replyMsg(writer, 400, err.Error(), err.Error())
	fmt.Println("Unmarshal:", err)
    }
}

func (s *SideCarProxy) StartServer() {
    r := mux.NewRouter()
    r.HandleFunc("/allocate/tokens", s.AllocateHandler)
    r.HandleFunc("/transfer", s.TransferHandler)
    r.HandleFunc("/balance", s.BalanceHandler)
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

func (s *SideCarProxy) SendEventToZrServer(from string, to string, value string) {
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
        /*resp, _ := */http.Post(s.HostName, body_type, req)
        //body, _ := ioutil.ReadAll(resp.Body)
        //fmt.Println(string(body))
    } else {
        fmt.Println(err)
    }
}

/*************************************************************************/



func init() {
   }

func balanceOf(address string) *big.Int {

    fmt.Println("call balance")
    balance, err := chain.Contract.BalanceOf(&bind.CallOpts{}, common.HexToAddress(address))
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(balance) // "1.0"

    return  balance
}

