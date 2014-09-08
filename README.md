a stupidly simple centralized peer list for groupcache

## Usage
### Server
Installing: `go get github.com/guregu/baka/bakad`

Running: `bakad [--bind 0.0.0.0:1337] [--timeout 1m]`

Timeout specifies how long to wait before a peer is considered dead and removed from the list. Uses the syntax specified [here](https://godoc.org/time#ParseDuration). 

## Client
```go
import 	"github.com/guregu/baka"

func InitCache(server, self string) {
	peers := groupcache.NewHTTPPool(self)
	// update the peer list every 10 seconds
	baka.Update(server, self, peers, time.Second*10)
}
```
`server` and `self` should be full URLs with port numbers. 