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
	stream proto.ChitChat_GetChitsServer
	id     string
	active bool
	error  chan error
	mu     sync.Mutex
}

type Pool struct {
	proto.UnimplementedChitChatServer
	Connection []*Chit_service
}

func main() {
	server := &Chit_service{}
	server.start_server()
}

func (p *Pool) CreateStream(pconn *proto.Connect, stream proto.ChitChat_GetChitsServer) error {
	conn := &Chit_service{
		stream: stream,
		id:     pconn.Author.Id,
		active: true,
		error:  make(chan error),
	}

	p.Connection = append(p.Connection, conn)

	return <-conn.error
}

func (s *Pool) BroadcastMessage(ctx context.Context, msg *proto.Chits) (*proto.Empty, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		wait.Add(1)

		go func(msg *proto.Chits, conn *Chit_service) {
			defer wait.Done()

			conn.mu.Lock()
			defer conn.mu.Unlock()

			if conn.active {
				err := conn.stream.Send(msg)
				log.Printf("Sending message to: %v from %v", conn.id, msg.Author)

				if err != nil {
					log.Printf("Error with Stream: %v - Error: %v\n", conn.stream, err)
					conn.active = false
					conn.error <- err
				}

			}
		}(msg, conn)

	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &proto.Empty{}, nil
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

func (server *Chit_service) SendChits(ctx context.Context, in *proto.Empty) {

} */

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
