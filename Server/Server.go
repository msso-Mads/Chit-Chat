package main

import (
	proto "Chit-Chat/gRPC"
	"context"
	"log"
	"net"
	"google.golang.org/grpc"
	"sync"
)

// Broadcaster is the struct which encompasses broadcasting
type Chit_service struct {
	proto.UnimplementedChitChatServer
	error  chan error
	broadcast chan string
	chatters map[string]chan *proto.Chits
	mu sync.Mutex
}



func main() {
	server := &Chit_service{}
	server.start_server()
}



func (server *Chit_service) JoinChit( in *proto.JoinRequest, 
	stream proto.ChitChat_JoinChitServer) error  {
	author := in.Author

	msgChan := make(chan *proto.Chits)

	server.mu.Lock()
	server.chatters[author.Name] = msgChan
	server.mu.Unlock()

	time := "your mom"
	log.Println("Participant", author.Name, "joined Chit Chat at logical time", time)
	joinMsg := &proto.Chits{
		Chit : "User has joined",
		Author : author.Name,
		TimeFormated: time,

	}

	if err := stream.Send(joinMsg); err != nil{
		log.Println(err)
		return err
	}

	for{ select{
	case msg := <- msgChan:
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
	
	log.Println("Participant", author.Name, "left Chit Chat at logical time your mom")

	ctx.Done();

	
	return &proto.Empty{

	}, nil
}

func (server *Chit_service) SendChits(ctx context.Context, in *proto.Chits) (*proto.Empty, error) {
	chit := in.Chit
	author := in.Author

	

	log.Println(author, ":", chit)
	for key, value := range server.chatters{
			if key != author{
				value <- in
			}
	}
	return &proto.Empty{
	}, nil
} 

func (server *Chit_service) start_server() {
	server.broadcast = make(chan string)
	server.chatters = make(map[string]chan *proto.Chits)
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
