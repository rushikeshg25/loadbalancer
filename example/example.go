package example

func main(){
	servers:=[]BackendServer{{
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
	lb:=Loadbalancer{
		Servers: servers,
		Strategy: ROUNDROBIN,
		Port: 9090,
		Host: "localhost",
		TotalRequests: 0,
	}
	lb.InitLoadbalancer();
}
