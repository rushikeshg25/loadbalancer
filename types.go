package loadbalancer

type BackendServer struct{
	Host string
	Port int
	IsHealthy bool
	ServerRequests int
}

type Loadbalancer struct{
	Servers []BackendServer
	Strategy StrategyType
	Port int
	Host string
	TotalRequests int
}

type StrategyType string

const (
	ROUNDROBIN StrategyType = "RoundRobin"
	LEASTCONN StrategyType = "LeastConn"
)

