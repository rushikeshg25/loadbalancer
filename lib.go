package loadbalancer

import (
	"log"
	"net"
)

func (lb *Loadbalancer) InitLoadbalancer(){
	listener,err:=net.Listen("tcp",string(lb.Port));
	if err!=nil{
		log.Fatalf("Error: %v",err)
	}

	for{
		conn,err:=listener.Accept()
		if err!=nil{
			log.Fatalf("Error: %v",err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn){
	defer conn.Close()
	server:=lb.Strategy.FindServer()

}