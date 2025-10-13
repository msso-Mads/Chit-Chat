package main

import (
	proto "Chit-Chat/gRPC"
	"context"
    "log"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main(){
    conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working client 1")
	}

	client := proto.NewChitChatClient(conn)
	
		
	join, err := client.JoinChit(context.Background(), 
		&proto.JoinRequest{
			Author: &proto.Author{
				Id: "1",
				Name: "Vee",
			},
		})
	if err != nil {
		log.Fatalf("Not working client 2")
	}

	log.Println(join)
	
}