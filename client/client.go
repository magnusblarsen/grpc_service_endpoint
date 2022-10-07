package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	gRPC "github.com/magnusblarsen/grpc_service_endpoint/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clientsName = flag.String("name", "default", "Senders name")
var serverPort = flag.String("server", "5400", "Tcp server")

var server gRPC.MyServiceClient //the server
var ServerConn *grpc.ClientConn //the server connection

func clientMain() {
	flag.Parse()

	fmt.Println("--- CLIENT APP ---")

	fmt.Println("--- join Server ---")
	ConnectToServer()
	defer ServerConn.Close()
}
func ConnectToServer() {
	var opts []grpc.DialOption
	opts = append(
		opts, grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	log.Printf("client %s: Attempts to dial on port %s\n", *clientsName, *serverPort)
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", *serverPort), opts...)
	if err != nil {
		log.Printf("Fail to Dial : %v", err)
		return
	}

	// makes a client from the server connection and saves the connection
	// and prints whether or not the connection is READY
	server = gRPC.NewMyServiceClient(conn)
	ServerConn = conn
	log.Println("the connection is: ", conn.GetState().String())
}

func GetTime() {
	info := &gRPC.Info{
		Clientname: *clientsName,
		Message:    "Hi this is a message from the client",
	}

	time, err := server.TellTime(context.Background(), info)
	if err != nil {
		log.Printf("Client %s: no response from the server, attempting to reconnect", *clientsName)
		log.Println(err)
	}
	fmt.Printf("The time received was %s", time)
}
