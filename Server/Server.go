package main

import (
	proto "Chit-Chat/gRPC"
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

// Broadcaster is the struct which encompasses broadcasting
type Chit_service struct {
	proto.UnimplementedChitChatServer
	error     chan error
	broadcast chan *proto.Chits
	chatters  map[string]chan *proto.Chits
	mu        sync.Mutex
}

func main() {
	server := &Chit_service{
		broadcast: make(chan *proto.Chits, 10),
		chatters:  make(map[string]chan *proto.Chits),
	}
	server.start_server()
}

func (server *Chit_service) JoinChit(in *proto.JoinRequest,
	stream proto.ChitChat_JoinChitServer) error {
	author := in.Author

	msgChan := make(chan *proto.Chits)

	server.mu.Lock()
	server.chatters[author.Name] = msgChan
	server.mu.Unlock()

	time := "your mom"
	log.Println("Participant", author.Name, "joined Chit Chat at logical time", time)

	msg := "Participant "
	msg += author.Name
	msg += " joined Chit Chat at logical time your mom"

	joinMsg := &proto.Chits{
		Chit:         msg,
		Author:       "system",
		TimeFormated: time,
	}

	server.broadcast <- joinMsg

	for {
		select {
		case msg := <-msgChan:
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

	msg := "Participant "
	msg += author.Name
	msg += " left Chit Chat at logical time your mom"

	chit := &proto.Chits{
		Chit:         msg,
		Author:       "system",
		TimeFormated: "your mom",
	}

	log.Println(msg)

	ctx.Done()

	server.broadcast <- chit

	return &proto.Empty{}, nil
}

func (server *Chit_service) SendChits(ctx context.Context, in *proto.Chits) (*proto.Empty, error) {
	chit := in.Chit
	author := in.Author

	log.Println(author, ":", chit)
	for key, value := range server.chatters {
		if key != author {
			value <- in
		}
	}
	return &proto.Empty{}, nil
}

func (server *Chit_service) start_server() {
	server.StartBroadcaster()
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")

	if err != nil {
		log.Fatalf("Did not work 1")
	}

	proto.RegisterChitChatServer(grpcServer, server)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work 2")
	}

	log.Println("the server has started")

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
