# tcp_proxy

Proxy and remote services into you applicaton to be to test how your application
reacts to network partitions of external/remote services. 

## Example Usage

```
import proxy "github.com/bubunyo/tcp_proxy"


func main(){
    // proxy a redis server running on port 6379 on localhost.
    p := proxy.New("localhost:6379") 

    // use the proxy's new address
    rdb := redis.NewClient(&redis.Options{Addr: p.Addr()})

    // test new connection
    err := rdb.Ping(context.Background()).Err()
    if err != nil {
        log.Panic("redis connection not established")
    }
}
```

## Todo

- Add tests.
