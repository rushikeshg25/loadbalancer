# TCP load balancer v1

`Loadbalancer.Serve(ctx, listener)` accepts TCP connections and proxies complete byte streams in both directions. Configure `Servers` with host, port and `IsHealthy: true`, and choose `ROUNDROBIN` or `LEASTCONN`. Both use shared reservations so simultaneous connections observe active load. Failed backend dials release reservations and try another eligible server before closing the client.

```go
lb := &loadbalancer.Loadbalancer{
    Strategy: loadbalancer.LEASTCONN,
    Servers: []loadbalancer.BackendServer{{Host: "127.0.0.1", Port: 8080, IsHealthy: true}},
}
listener, err := net.Listen("tcp", "127.0.0.1:9090")
if err != nil { panic(err) }
err = lb.Serve(ctx, listener)
```

Cancelling the context closes the listener and active connections, then waits for proxy goroutines. TCP half-closes are propagated so peers can finish their response. Backend dialing has a three-second timeout. Stream duration is otherwise unrestricted until the context is cancelled. No request buffering or 1 KiB truncation is used.

Configure before calling Serve; do not modify/read mutable server counters while serving. Use one Serve call per instance. `IsHealthy` is static eligibility, not an active health probe. TLS termination, dynamic discovery, live configuration and connection admission limits are outside v1. The legacy `InitLoadbalancer` convenience method blocks without a cancellation handle; new applications should use Serve and handle its errors.

Run `go test -race ./...` to exercise real localhost streaming and cancellation. `go run ./example` starts the legacy sample on port 9090.
