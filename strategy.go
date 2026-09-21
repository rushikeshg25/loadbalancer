package loadbalancer

// NextServer selects a healthy backend without reserving a live connection.
func (lb *Loadbalancer) NextServer() BackendServer {
	return lb.selectServer(lb.Strategy)
}
func (lb *Loadbalancer) selectServer(strategy StrategyType) BackendServer {
	i := lb.reserveWith(strategy, nil)
	if i < 0 {
		return BackendServer{}
	}
	lb.mu.Lock()
	s := lb.Servers[i]
	lb.mu.Unlock()
	lb.release(i)
	return s
}
func (lb *Loadbalancer) RoundRobin() BackendServer { return lb.selectServer(ROUNDROBIN) }
func (lb *Loadbalancer) LeastConn() BackendServer  { return lb.selectServer(LEASTCONN) }
