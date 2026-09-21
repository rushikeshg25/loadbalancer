package loadbalancer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Serve owns listener until ctx is cancelled. Configuration must not change while serving.
func (lb *Loadbalancer) Serve(ctx context.Context, listener net.Listener) error {
	if len(lb.Servers) == 0 {
		return errors.New("no backends configured")
	}
	for _, s := range lb.Servers {
		if s.Host == "" || s.Port < 1 || s.Port > 65535 {
			return errors.New("invalid backend address")
		}
	}
	if lb.Strategy != "" && lb.Strategy != ROUNDROBIN && lb.Strategy != LEASTCONN {
		return errors.New("unknown balancing strategy")
	}
	defer listener.Close()
	stop := context.AfterFunc(ctx, func() { listener.Close() })
	defer stop()
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	defer func() { cancel(); wg.Wait() }()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		wg.Add(1)
		go func() { defer wg.Done(); lb.proxy(ctx, conn) }()
	}
}
func (lb *Loadbalancer) InitLoadbalancer() {
	listener, err := net.Listen("tcp", net.JoinHostPort(lb.Host, fmt.Sprint(lb.Port)))
	if err != nil {
		return
	}
	_ = lb.Serve(context.Background(), listener)
}
func (lb *Loadbalancer) reserve(excluded map[int]bool) int {
	return lb.reserveWith(lb.Strategy, excluded)
}
func (lb *Loadbalancer) reserveWith(strategy StrategyType, excluded map[int]bool) int {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	best := -1
	for n := 0; n < len(lb.Servers); n++ {
		i := (lb.TotalRequests + n) % len(lb.Servers)
		s := lb.Servers[i]
		if excluded[i] || !s.IsHealthy {
			continue
		}
		if best < 0 || strategy == LEASTCONN && s.ServerRequests < lb.Servers[best].ServerRequests {
			best = i
		}
		if strategy != LEASTCONN {
			break
		}
	}
	if best >= 0 {
		lb.Servers[best].ServerRequests++
		lb.TotalRequests = (best + 1) % len(lb.Servers)
	}
	return best
}
func (lb *Loadbalancer) release(i int) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.Servers[i].ServerRequests--
}
func (lb *Loadbalancer) proxy(ctx context.Context, client net.Conn) {
	defer client.Close()
	stopClient := context.AfterFunc(ctx, func() { client.Close() })
	defer stopClient()
	excluded := map[int]bool{}
	for len(excluded) < len(lb.Servers) {
		i := lb.reserve(excluded)
		if i < 0 {
			return
		}
		excluded[i] = true
		lb.mu.Lock()
		s := lb.Servers[i]
		lb.mu.Unlock()
		backend, err := (&net.Dialer{Timeout: 3 * time.Second}).DialContext(ctx, "tcp", net.JoinHostPort(s.Host, fmt.Sprint(s.Port)))
		if err != nil {
			lb.release(i)
			continue
		}
		defer lb.release(i)
		defer backend.Close()
		stopBackend := context.AfterFunc(ctx, func() { backend.Close() })
		defer stopBackend()
		done := make(chan struct{})
		go func() {
			_, err := io.Copy(backend, client)
			if err != nil {
				backend.Close()
				client.Close()
			} else {
				closeWrite(backend)
			}
			close(done)
		}()
		_, err = io.Copy(client, backend)
		if err != nil {
			backend.Close()
			client.Close()
		} else {
			closeWrite(client)
		}
		<-done
		return
	}
}
func closeWrite(c net.Conn) {
	if c, ok := c.(interface{ CloseWrite() error }); ok {
		_ = c.CloseWrite()
	}
}
