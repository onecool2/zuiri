package main
import  (
	"github.com/onecool2/zuiri/sidecar/pkg/sidecarproxy"
	"flag"
	"github.com/golang/glog"
	"github.com/onecool2/zuiri/sidecar/pkg/chain"
)

const (
	ZUI_RI_SERVER_HOST string = "http://127.0.0.1:3001/event"
	EVENT string = "event"
	BALANCE string = "balance"
	ALLOCATETOKENS string = "allocatetokens"
	TRANSFER string = "transfer"
)

func main() {
    //init()
    flag.Parse()
    defer glog.Flush()

    sideCarProxy := sidecarproxy.SideCarProxy{ZUI_RI_SERVER_HOST, EVENT, BALANCE, ALLOCATETOKENS, TRANSFER}
    go sideCarProxy.StartServer()
    go chain.LoopAndSendTx()
    go chain.GoThroughBlock()
    select {}
}
