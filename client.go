package main

import (
	"bufio"
	"fmt"
	"log"
	rpc "net/rpc"
	"os"
	"strings"
)

type RegisterArgs struct {
	UserName string
}

type RegisterReply struct {
	Success bool
	ID      int
}

type SendMessageArgs struct {
	UserName string
	Content  string
}

func main() {
	client, err := rpc.Dial("tcp", "0.0.0.0:42586")
	if err != nil {
		log.Fatal("Error connecting:", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter your name: ")
	scanner.Scan()
	name := strings.TrimSpace(scanner.Text())
	if name == "" {
		name = "Anonymous"
	}

	// Register on server
	var regReply RegisterReply
	err = client.Call("ChatServer.Register", &RegisterArgs{UserName: name}, &regReply)
	if err != nil || !regReply.Success {
		log.Fatal("Failed to register with server")
	}

	fmt.Printf("Connected as User [%d] (%s)\n", regReply.ID, name)

	// Receiving goroutine
	go func() {
		for {
			var msg string
			err := client.Call("ChatServer.Receive", name, &msg)
			if err != nil {
				fmt.Println("Disconnected:", err)
				return
			}
			fmt.Println(msg)
		}
	}()

	// Sending loop
	for {
		fmt.Print("> ")
		scanner.Scan()
		text := scanner.Text()

		if text == "exit" {
			fmt.Println("Goodbye!")
			return
		}

		args := SendMessageArgs{
			UserName: name,
			Content:  text,
		}

		var ok bool
		err := client.Call("ChatServer.SendMessage", &args, &ok)
		if err != nil {
			fmt.Println("Send failed:", err)
		}
	}
}
