package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main(){
	listener,err:= net.Listen("tcp",":9000")
	if err!=nil{
	log.Fatal(err)
	}
	defer listener.Close()
	log.Println("Server is listening on port 9000...")

	for{
		conn,err:=listener.Accept()
		if err!=nil{
			fmt.Println("Connection Error")
			continue
		}
		go handleClient(conn)
	}

}

func handleClient(conn net.Conn){
	defer conn.Close()
	fmt.Println("New client connected:",conn.RemoteAddr())

	scanner:=bufio.NewReader(conn)
	for{
		line,err:=scanner.ReadString('\n')
		if err!=nil{
			fmt.Println("Client disconnected:",err)
			break
		}
		fmt.Println("Received from client:",line)
		response:=processCommand(line)
		conn.Write([]byte(response+'\n'))

	}
}