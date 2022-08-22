# tcp_proxy

Proxy and remote services into you applicaton to be to test how your application
reacts to network partitions of external/remote services. 

## Example Usage

```
import proxy "github.com/bubunyo/tcp_proxy"


func main(){
    p := proxy.New("localhost:6379") // proxy a redis server running on port 6379 on localhost.

    rdb := redis.NewClient(&redis.Options{Addr: p.Addr()}) // use the proxy's new address

    if err:=rdb.Ping(context.Background()).Err(); err != nil {
        log.Panic("redis connection not established")
    }
```

## Todo

- Add tests.
