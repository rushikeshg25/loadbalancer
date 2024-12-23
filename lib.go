package loadbalancer

import (
	"fmt"
	"log"
	"net"
)

func (lb *Loadbalancer) InitLoadbalancer(){
	listener,err:=net.Listen("tcp",fmt.Sprintf(":%d",lb.Port));
	if err!=nil{
		log.Fatalf("Error: %v",err)
	}
	log.Printf("Loadbalancer started on %s:%d",lb.Host,lb.Port)
	for{
		conn,err:=listener.Accept()
		if err!=nil{
			log.Fatalf("Error: %v",err)
		}
		nextServer:=lb.NextServer()
		log.Printf("Request sent to %s:%d",nextServer.Host,nextServer.Port)
		go handleConnection(conn,&nextServer)

	}
}

func handleConnection(conn net.Conn,server *BackendServer){
	bytes:=make([]byte,1024)
	conn.Read(bytes)
	defer conn.Close()
	server.ServerRequests++;

	c,err:=net.Dial("tcp",fmt.Sprintf("%s:%d",server.Host,server.Port))
	if err!=nil{
		log.Fatalf("Error: %v",err)
	}
	defer c.Close()
	_,err=c.Write(bytes)
	if err!=nil{
		log.Fatalf("Error: %v",err)
		conn.Write([]byte("Error with Backend Server"))
		conn.Close()
		return
	}
	c.Read(bytes);
	// log.Println(string(bytes))
	conn.Write(bytes);
}