package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// 1. Try to dial with a 5-second timeout
	address := "127.0.0.1:9000"
	fmt.Printf("Attempting to connect to %s...\n", address)
	
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		fmt.Printf("❌ CONNECTION FAILED: %v\n", err)
		fmt.Println("Check if your server terminal shows 'Server is listening...'")
		return
	}
	defer conn.Close()
	fmt.Println("✅ CONNECTED! Type your commands (e.g., SET key value)")

	// Read from Server, Print to Terminal
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println("\nServer:", scanner.Text())
			fmt.Print("> ") // Print prompt for user
		}
	}()

	// Read from Terminal, Send to Server
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			return
		}
		fmt.Fprintf(conn, text+"\n")
	}
}