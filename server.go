package main

import (
	"fmt"
	"log"
	"net"
	rpc "net/rpc"
	"sync"
	"time"
)

type ChatServer struct{}

// Client struct with numeric ID
type Client struct {
	ID   int
	Name string
	Ch   chan string
}

var (
	clients   = make(map[string]*Client) // key = username
	clientMux sync.Mutex

	broadcast = make(chan string)

	nextID = 1 // auto-increment ID counter
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

func (s *ChatServer) Register(args *RegisterArgs, reply *RegisterReply) error {
	clientMux.Lock()

	id := nextID
	nextID++

	clients[args.UserName] = &Client{
		ID:   id,
		Name: args.UserName,
		Ch:   make(chan string, 10),
	}

	clientMux.Unlock()

	// Notify all other clients
	broadcast <- fmt.Sprintf("User [%d] (%s) joined", id, args.UserName)

	reply.Success = true
	reply.ID = id
	return nil
}

func (s *ChatServer) SendMessage(args *SendMessageArgs, reply *bool) error {
	clientMux.Lock()
	c, ok := clients[args.UserName]
	clientMux.Unlock()

	if !ok {
		*reply = false
		return fmt.Errorf("user not registered")
	}

	msg := fmt.Sprintf("[%s] User [%d] (%s): %s",
		time.Now().Format("15:04:05"),
		c.ID,
		c.Name,
		args.Content,
	)

	// Send to everyone except sender
	clientMux.Lock()
	for name, cl := range clients {
		if name == args.UserName {
			continue
		}
		cl.Ch <- msg
	}
	clientMux.Unlock()

	*reply = true
	return nil
}

func (s *ChatServer) Receive(userName string, msg *string) error {
	clientMux.Lock()
	c, ok := clients[userName]
	clientMux.Unlock()

	if !ok {
		return fmt.Errorf("client not registered")
	}

	*msg = <-c.Ch
	return nil
}

func main() {
	// Global broadcast goroutine
	go func() {
		for msg := range broadcast {
			clientMux.Lock()
			for _, c := range clients {
				c.Ch <- msg
			}
			clientMux.Unlock()
		}
	}()

	server := new(ChatServer)
	rpc.Register(server)

	addr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:42586")
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Real-time chat server running on port 42586...")
	rpc.Accept(listener)
}
