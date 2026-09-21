package loadbalancer

import "sync"

type BackendServer struct {
	Host           string
	Port           int
	IsHealthy      bool
	ServerRequests int
}

type Loadbalancer struct {
	mu            sync.Mutex
	Servers       []BackendServer
	Strategy      StrategyType
	Port          int
	Host          string
	TotalRequests int
}

type StrategyType string

const (
	ROUNDROBIN StrategyType = "RoundRobin"
	LEASTCONN  StrategyType = "LeastConn"
)
