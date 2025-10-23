package main

import (
	proto "Chit-Chat/gRPC"
	"bufio"
	"context"
	"log"
	"os"
	"os/user"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"
)

func main() {
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
				Id:   uid,
				Name: username,
			},
		})
	if err != nil {
		log.Fatalf("Not working client 2")
	}


	go recieve(join)

	send(client, username)

}

func send(client proto.ChitChatClient, username string) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		s := scanner.Text()
		if s == "Leave" {
			break
		} else {
			send, err := client.SendChits(context.Background(),
				&proto.Chit{
					Chit:         s,
					Author:       username,
					TimeFormated: time.Now().String(),
				},
			)
			if err != nil {
				log.Fatalf("client not sending message")
			}

			log.Println(send)

		}
	}
}

func recieve( stream grpc.ServerStreamingClient[proto.Chits]){
	for  {
		response, err := stream.Recv()
		if err == io.EOF{
			break
		}
		if err != nil{
			log.Fatalf("error while recieving message")
		}
		fmt.Println(response)
	}
}