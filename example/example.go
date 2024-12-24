package main

import (
	"github.com/rushikeshg25/loadbalancer"
)

func main(){
	servers:=[]loadbalancer.BackendServer{{
		Host: "localhost",
		Port: 8080,
		IsHealthy: true,
		ServerRequests: 0,
	},
	{
		Host: "localhost",
		Port: 8081,
		IsHealthy: true,
		ServerRequests: 0,
	}}
	lb:=loadbalancer.Loadbalancer{
		Servers: servers,
		Strategy: loadbalancer.ROUNDROBIN,
		Port: 9090,
		Host: "localhost",
		TotalRequests: 0,
	}
	lb.InitLoadbalancer();
}
