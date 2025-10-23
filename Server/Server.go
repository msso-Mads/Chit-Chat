package main

import (
	proto "Chit-Chat/gRPC"
	"context"
	"log"
	"net"
	"sync"
	"time"
	"google.golang.org/grpc"
)

// Broadcaster is the struct which encompasses broadcasting
type Chit_service struct {
	proto.UnimplementedChitChatServer
	stream proto.ChitChat_GetChitsServer
	id     string
	active bool
	error  chan error
	mu     sync.Mutex
	chits []string
	author []string
	time []string
}



func main() {
	server := &Chit_service{}
	server.start_server()
}



func (server *Chit_service) JoinChit(ctx context.Context, in *proto.JoinRequest) (*proto.Join, error) {
	author := in.Author
	time := "your mom"
	log.Println("Participant", author, "joined Chit Chat at logical time", time)
	return &proto.Join{
		Author: author,
		Time:   time,
	}, nil
}

/*
func (server *Chit_service) LeaveChit(ctx context.Context, in *proto.Empty) {

}

func (server *Chit_service) GetChits(ctx context.Context, in *proto.Empty){

}
*/
func (server *Chit_service) SendChits(ctx context.Context, in *proto.Chit) (*proto.Empty, error) {
	chit := in.Chit
	author := in.Author
	server.chits = append(server.chits, chit)
	server.author = append(server.author, author)
	server.time = append(server.time, time.Now().String())
	log.Println(author, ":", chit)
	return &proto.Empty{
	}, nil
} 

func (server *Chit_service) start_server() {
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
