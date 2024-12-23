package loadbalancer

type BackendServer struct{
	Host string
	Port int
	IsHealthy bool
	Requests int
}

type Loadbalancer struct{
	Servers []BackendServer
	Strategy Strategy
	Port int
	Host string
}

type StrategyType string

const (
	ROUNDROBIN StrategyType = "RoundRobin"
	LEASTCONN StrategyType = "LeastConn"
)

type Strategy struct{
	Type StrategyType
}