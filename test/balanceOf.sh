#curl -XPOST "http://127.0.0.1:3000/balance" -H 'Content-Type: application/json' -d'
curl -XPOST "http://127.0.0.1:3000/balance" -w "http_code: %{http_code}  content_type:%{content_type} \n" -H 'application/x-www-form-urlencoded' -d'
{
	"address": "0x2d9c95dd961ea36350f95ed812669da85897a9d8"
}'
#curl -XPOST "http://127.0.0.1:3001/event" -H 'Content-Type: application/json' -d'
#{
#	"from": "update",
#	"to": "update",
#	"value": "333",
#	"time": "333"
#}'
