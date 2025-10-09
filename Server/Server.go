package main

import (
	proto "Chit-Chat/gRPC"
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type Chit_service struct {
	proto.UnimplementedChitChatServer
}

// Broadcaster is the struct which encompasses broadcasting
type Broadcaster struct {
	cond        *sync.Cond
	subscribers map[interface{}]func(interface{})
	message     interface{}
	running     bool
}

func main() {
	server := &Chit_service{}

	server.start_server()
}

func (server *Chit_service) join_server(ctx context.Context, in *proto.Empty) (*proto.Join, error) {

}

func (server *Chit_service) leave_server(ctx context.Context, in *proto.Empty) {

}

func (server *Chit_service) send_chit()

func (server *Chit_service) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", "number of localhost")

	if err != nil {
		log.Fatalf("Did not work")
	}

	proto.RegisterChitChatServer(grpcServer, server)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}
}

// SetupBroadcaster gives the broadcaster object to be used further in messaging
func SetupBroadcaster() *Broadcaster {

	return &Broadcaster{
		cond:        sync.NewCond(&sync.RWMutex{}),
		subscribers: map[interface{}]func(interface{}){},
	}
}
