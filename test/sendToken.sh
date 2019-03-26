curl -XPOST "http://127.0.0.1:3000/transfer" -H 'Content-Type: application/json' -d'
{
	"from": "0x2d9c95dd961ea36350f95ed812669da85897a9d8",	    
	"to": "0x1F959fe299ed89479fb250C6bBDd7d8ef9bdd501",
	"value": "1",
	"serialNumber": "333"
}'
curl -XPOST "http://127.0.0.1:3000/transfer" -H 'Content-Type: application/json' -d'
{
	"from": "0x1F959fe299ed89479fb250C6bBDd7d8ef9bdd501",
	"to": "0x2d9c95dd961ea36350f95ed812669da85897a9d8",
	"value": "1",
	"serialNumber": "333"
}'
    #sleep 1

#curl -XPOST "http://127.0.0.1:3000/balance" -H 'Content-Type: application/json' -d'
#{
#	"address": "0x96216849c49358B10257cb55b28eA603c874b05E"
#}'

#curl -XPOST "http://127.0.0.1:3001/event" -H 'Content-Type: application/json' -d'
#{
#	"from": "update",
#	"to": "update",
#	"value": "333",
#	"time": "333"
#}'
