package main

import (
	proto "Chit-Chat/gRPC"
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// Broadcaster is the struct which encompasses broadcasting
type Chit_service struct {
	proto.UnimplementedChitChatServer
	error       chan error
	broadcast   chan *proto.Chits
	chatters    map[string]chan *proto.Chits
	mu          sync.Mutex
	logicalTime int64
	grpc        *grpc.Server
}

func main() {
	server := &Chit_service{
		broadcast:   make(chan *proto.Chits, 10),
		chatters:    make(map[string]chan *proto.Chits),
		logicalTime: 0,
	}
	server.start_server()
}

func (server *Chit_service) JoinChit(in *proto.JoinRequest,
	stream proto.ChitChat_JoinChitServer) error {
	server.logicalTime = max(server.logicalTime, in.Time) + 1
	author := in.Author

	msgChan := make(chan *proto.Chits)

	server.mu.Lock()
	server.chatters[author.Name] = msgChan
	server.mu.Unlock()

	log.Println("Participant", author.Name, "joined Chit Chat at logical time", server.logicalTime)

	msg := "Participant " + author.Name + " joined Chit Chat at logical time " + strconv.FormatInt(server.logicalTime, 10)

	joinMsg := &proto.Chits{
		Chit:         msg,
		Author:       "Server",
		TimeFormated: server.logicalTime,
	}

	server.broadcast <- joinMsg

	for {
		select {

		case msg := <-msgChan:
			server.logicalTime++
			msg.TimeFormated = server.logicalTime
			if err := stream.Send(msg); err != nil {
				return err
			}

		case <-stream.Context().Done():
			log.Println("Stream ended")
			return nil

		}

	}

}

func (server *Chit_service) LeaveChit(ctx context.Context, in *proto.Leave) (*proto.Empty, error) {
	author := in.Author

	server.logicalTime = max(server.logicalTime, in.Time) + 1

	msg := "Participant " + author.Name + " left Chit Chat at logical time " + strconv.FormatInt(server.logicalTime, 10)

	chit := &proto.Chits{
		Chit:         msg,
		Author:       "Server",
		TimeFormated: server.logicalTime,
	}

	log.Println(msg)

	ctx.Done()
	delete(server.chatters, author.Name)

	server.broadcast <- chit
	
	if len(server.chatters) == 0 {
		log.Println("Server is closing")
		go func(){
			time.Sleep(10*time.Second)
			log.Println("Server stopped")
			server.grpc.Stop()
			time.Sleep(10*time.Second)
			os.Exit(0)
		}()
		
		
	}

	

	return &proto.Empty{}, nil
}

func (server *Chit_service) SendChits(ctx context.Context, in *proto.Chits) (*proto.Empty, error) {
	chit := in.Chit
	author := in.Author

	server.logicalTime = max(server.logicalTime, in.TimeFormated) + 1

	log.Println("Server:",author, "sent chit:", chit, " - Logical time", server.logicalTime)

	for _, value := range server.chatters {
		value <- in

	}
	return &proto.Empty{}, nil
}

func (server *Chit_service) start_server() {
	server.StartBroadcaster()
	server.grpc = grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")

	if err != nil {
		log.Fatalf("Did not work 1")
	}

	log.Println("the server has started")

	proto.RegisterChitChatServer(server.grpc, server)

	err = server.grpc.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work 2")
	}

}

func (server *Chit_service) StartBroadcaster() {
	go func() {
		log.Println("broadcaster started")
		for msg := range server.broadcast {
			server.mu.Lock()
			for _, ch := range server.chatters {
				select {
				case ch <- msg:
				default:
					// Avoid blocking slow clients
				}

			}
			server.mu.Unlock()
		}
	}()
}
