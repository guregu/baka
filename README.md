a stupidly simple centralized peer list for groupcache

## Usage
### Server
Installing: `go get github.com/guregu/baka/bakad`

Running: `bakad [--bind 0.0.0.0:1337] [--timeout 1m]`

Timeout specifies how long to wait before a peer is considered dead and removed from the list. Uses the syntax specified [here](https://godoc.org/time#ParseDuration). 

## Client
```go
import 	"github.com/guregu/baka"

var peers *groupcache.HTTPPool

func main() {
	// e.x. bakad is at 10.0.0.1:1337 and this server is at 10.0.0.42:7000
	server := "http://10.0.0.1:1337"
	self := "http://10.0.0.42:7000"
	peers = groupcache.NewHTTPPool(self)
	// refresh peers list every 10 seconds
	baka.Update(server, self, peers, time.Second*10)
}
```