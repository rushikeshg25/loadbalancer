package loadbalancer

// NextServer selects a healthy backend without reserving a live connection.
func (lb *Loadbalancer) NextServer() BackendServer {
	i := lb.reserve(nil)
	if i < 0 {
		return BackendServer{}
	}
	lb.mu.Lock()
	s := lb.Servers[i]
	lb.mu.Unlock()
	lb.release(i)
	return s
}
func (lb *Loadbalancer) RoundRobin() BackendServer { return lb.NextServer() }
func (lb *Loadbalancer) LeastConn() BackendServer  { return lb.NextServer() }
