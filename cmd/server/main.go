package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/AbhinavG786/Gopher-Cache.git/internal/aof"
	store "github.com/AbhinavG786/Gopher-Cache.git/internal/engine"
	parser "github.com/AbhinavG786/Gopher-Cache.git/internal/protocol"
	"github.com/AbhinavG786/Gopher-Cache.git/internal/pubsub"
)

func main(){
	log.Println("Step 1: Engine Init")
	store:=store.New()
	log.Println("Hub Init")
	hub:=pubsub.NewHub()
	log.Println("Step 2: AOF Init")
	aofFile,err:=aof.NewAOF("database.aof")
	if err!=nil{
		log.Fatal("Failed to open AOF file:",err)
	}
	defer aofFile.Close()
	log.Println("Restoring data from AOF")
	err=aofFile.Replay(func(line string) {
		parser.Process(store,hub, line)
	})
	if err!=nil{
		log.Println("AOF Replay finished with Error/EOF:",err)
	} else{
		log.Println("AOF Replay finished successfully")
	}
	log.Println("Step 4: Attempting to Listen on 0.0.0.0:9000")
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
		go handleClient(conn,store,aofFile,hub)
	}

}

func handleClient(conn net.Conn,store parser.StoreInterface,aofFile *aof.AOF,hub *pubsub.Hub){
	var subscribedTopic string
	var userChan chan string
	defer func(){
		if subscribedTopic!="" && userChan!=nil{
			hub.Unsubscribe(subscribedTopic,userChan)
		}
	conn.Close()
	}()
	fmt.Println("New client connected:",conn.RemoteAddr())

	scanner:=bufio.NewReader(conn)
	for{
		line,err:=scanner.ReadString('\n')
		if err!=nil{
			fmt.Println("Client disconnected:",err)
			break
		}
		line=strings.TrimSpace(line)
		if line==""{
			continue
		}
		fmt.Println("Received from client:",line)
		response, ch:=parser.Process(store,hub, line)
		cmdUpper:=strings.ToUpper(strings.Fields(line)[0])
		if cmdUpper=="SET" || cmdUpper=="DEL"{
			if response=="Key set" || response=="Key deleted"{
				err:=aofFile.Append(line)
				if err!=nil{
					log.Println("Error occurred while appending to AOF:",err)
				}
			}
		}
		if ch!=nil{
			parts:=strings.Fields(line)
			subscribedTopic=parts[1]
			userChan=ch
			go func(topic string,c chan string){
				for msg:=range c{
					fmt.Fprintf(conn, "\n[PUB %s]: %s\n> ", topic, msg)
				}
			}(subscribedTopic,userChan)
		}

		conn.Write([]byte(response+"\n"))

	}
}