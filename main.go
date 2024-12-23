package loadbalancer

import (
	"log"
)

func main(){
	log.Printf("")
	server1:=[]BackendServer{BackendServer{
		Host: "localhost",
		Port: 8080,
		IsHealthy: true,
		Requests: 0,
	}}

}