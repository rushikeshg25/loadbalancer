package loadbalancer

func (lb *Loadbalancer) RoundRobin() BackendServer {
	r:=lb.Servers[lb.TotalRequests%len(lb.Servers)]
	lb.TotalRequests++
	return r;
}
func (lb *Loadbalancer) LeastConn() BackendServer {
	return lb.Servers[0]
}


func (lb *Loadbalancer) NextServer() BackendServer {
	switch lb.Strategy {
	case ROUNDROBIN:
		return lb.RoundRobin()
	case LEASTCONN:
		return lb.LeastConn()
	default:
		return lb.RoundRobin()
	}
}

