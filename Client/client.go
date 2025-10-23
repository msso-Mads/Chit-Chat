package main

import (
	proto "Chit-Chat/gRPC"
	"context"
	"log"
	"os/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"fmt"
)

func main(){
    conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working client 1")
	}

    currentUser, err := user.Current()
    if err != nil {
        log.Fatalf(err.Error())
    }

	client := proto.NewChitChatClient(conn)
	username := currentUser.Username
	uid := currentUser.Uid
		
	join, err := client.JoinChit(context.Background(), 
		&proto.JoinRequest{
			Author: &proto.Author{
				Id: uid,
				Name: username,
			},
		})
	if err != nil {
		log.Fatalf("Not working client 2")
	}
	
	log.Println(join)
	send(client, username)
}

func send(client proto.ChitChatClient, username string){
	for{
		var command string
		fmt.Scanln(&command)
		if (command == "Leave"){
			break
		} else {
			send, err := client.SendChits(context.Background(),
			&proto.Chit{
				Chit: command,
				Author: username,
				TimeFormated: "your mom",
			},
		)
		if err != nil {
			log.Fatalf("client not sending message")
		}

		log.Println(send)
			
		}
	}
}