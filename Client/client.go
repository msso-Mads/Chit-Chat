package main

import (
	proto "Chit-Chat/gRPC"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logicalTime int64 = 0

func main() {
	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working client 1")
	}
	log.Println("Client has connected to server")
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	client := proto.NewChitChatClient(conn)
	fmt.Println("Enter your username")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := scanner.Text()
	username := s
	uid := currentUser.Uid

	logicalTime++
	join, err := client.JoinChit(context.Background(),
		&proto.JoinRequest{
			Author: &proto.Author{
				Id:   uid,
				Name: username,
			},
			Time: logicalTime,
		})
	if err != nil {
		log.Fatalf("Not working client 2")
	}

	go recieve(join)

	send(client, username, uid)

}

func send(client proto.ChitChatClient, username string, uid string) {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		s := scanner.Text()
		if s == "Leave" {
			break

		} else if len(s) > 128 {
			log.Println("Message is too long.")
		} else {
			logicalTime++
			send, err := client.SendChits(context.Background(),
				&proto.Chits{
					Chit:         s,
					Author:       username,
					TimeFormated: logicalTime,
				},
			)
			if err != nil {
				log.Fatalf("client not sending message")
			}

			log.Println(send)

		}
	}
	logicalTime++
	send, err := client.LeaveChit(context.Background(),
		&proto.Leave{
			Author: &proto.Author{
				Id:   uid,
				Name: username,
			},
			Time: logicalTime,
		},
	)
	if err != nil {
		log.Fatalf("client not leaving")
	}
	log.Println(send)
	log.Println("Client has left the server")
}

func recieve(stream grpc.ServerStreamingClient[proto.Chits]) {
	for {
		response, err := stream.Recv()
		logicalTime = max(logicalTime, response.TimeFormated) + 1
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while recieving message")
		}
		fmt.Println(response.Author, ":", response.Chit , " - Logical Time:" , logicalTime)
	}
}
