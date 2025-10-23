package main

import (
	proto "Chit-Chat/gRPC"
	"context"
	"log"
	"net"
	"time"
	"google.golang.org/grpc"
)

// Broadcaster is the struct which encompasses broadcasting
type Chit_service struct {
	proto.UnimplementedChitChatServer
	stream proto.ChitChat_JoinChitServer
	error  chan error
	chits []string
	author []string
	time []string
	broadcast chan string
}



func main() {
	server := &Chit_service{}
	server.start_server()
}



func (server *Chit_service) JoinChit(in *proto.JoinRequest, stream proto.ChitChat_JoinChitServer) error  {
	author := in.Author
	time := "your mom"
	log.Println("Participant", author, "joined Chit Chat at logical time", time)
	joinMsg := &proto.Chits{
		Chit : "User has joined",
		Author : author.Name,
		TimeFormated: time,

	}
	if err := stream.Send(joinMsg); err != nil{
		log.Println(err)
		return err
	}

	

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
