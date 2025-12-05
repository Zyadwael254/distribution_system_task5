# distribution_system_task5
RPC chat with server->client broadcasts.

How it works:

Each client starts a small RPC server to receive delivered messages (Client.Receive).
Client registers with central server (ChatServer.Register) giving its ID and RPC address.
Server keeps a map of connected clients, a broadcast channel and a goroutine that sends messages to all clients (except the sender).
On join: server appends "User [ID] joined" to history and broadcasts it to others (no self-echo).
On send: server appends "Sender: text" to history and broadcasts to other clients (no self-echo).
History RPC still available: ChatServer.History.
Run locally:

Create new folder, save server.go and client.go.

Start server: go run server.go

Start clients (each in its own terminal): go run client.go --name Alice go run client.go --name Bob

Client commands:

Type a message and press Enter to send.
Type history to fetch full chat history from server.
Type exit to quit.
