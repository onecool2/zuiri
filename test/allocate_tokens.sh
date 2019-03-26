curl -XPOST "http://127.0.0.1:3000/allocate/tokens" -H 'application/x-www-form-urlencoded' -d'
{
	"to": "0x2d9c95dd961ea36350f95ed812669da85897a9d8",
	"value": "1"
}'

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
