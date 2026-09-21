package loadbalancer

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"
	"time"
)

func TestStreamAndShutdown(t *testing.T) {
	backend, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()
	go func() {
		c, e := backend.Accept()
		if e == nil {
			defer c.Close()
			io.Copy(c, c)
		}
	}()
	port := backend.Addr().(*net.TCPAddr).Port
	lb := &Loadbalancer{Servers: []BackendServer{{Host: "127.0.0.1", Port: port, IsHealthy: true}}}
	front, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- lb.Serve(ctx, front) }()
	c, err := net.Dial("tcp", front.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	c.SetDeadline(time.Now().Add(3 * time.Second))
	payload := bytes.Repeat([]byte("message"), 20000)
	go c.Write(payload)
	got := make([]byte, len(payload))
	if _, err = io.ReadFull(c, got); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatal("corrupt stream")
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("shutdown hung")
	}
}
func TestStrategies(t *testing.T) {
	lb := &Loadbalancer{Strategy: LEASTCONN, Servers: []BackendServer{{IsHealthy: true}, {IsHealthy: true}}}
	a := lb.reserve(nil)
	b := lb.reserve(nil)
	if a == b {
		t.Fatal("least connections reused busy backend")
	}
	lb.release(a)
	if lb.reserve(nil) != a {
		t.Fatal("did not reuse idle backend")
	}
	lb.Servers[0].IsHealthy = false
	if lb.reserve(map[int]bool{1: true}) != -1 {
		t.Fatal("selected unhealthy backend")
	}
}
func TestEmptyConfiguration(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	if (&Loadbalancer{}).Serve(context.Background(), l) == nil {
		t.Fatal("accepted empty config")
	}
}
