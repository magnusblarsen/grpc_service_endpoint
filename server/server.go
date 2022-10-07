package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	gRPC "github.com/magnusblarsen/grpc_service_endpoint/proto"
	"google.golang.org/grpc"
)

type Server struct {
	// an interface that the server needs to have
	gRPC.UnimplementedMyServiceServer

	// here you can impliment other fields that you want
	name             string
	port             string
	numberOfRequests int
}

// flags are used to get arguments from the terminal. Flags take a value, a default value and a description of the flag.
// to use a flag then just add it as an argument when running the program.
var serverName = flag.String("name", "default", "Senders name") // set with "-name <name>" in terminal
var port = flag.String("port", "4500", "Server port")           // set with "-port <port>" in terminal

func serverMain() {
	flag.Parse()
	fmt.Println(".:server is starting:.")
	launchServer()
}

func launchServer() {
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *port))
	if err != nil {
		log.Printf("Server %s: Failed to listen on port %s: %v", *serverName, *port, err)
	}
	grpcServer := grpc.NewServer()

	server := &Server{
		name:             *serverName,
		port:             *port,
		numberOfRequests: 0,
	}

	gRPC.RegisterMyServiceServer(grpcServer, server)
	log.Printf("Server %s: Listening at %v\n", *serverName, list.Addr())
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

func (s *Server) TellTime(ctx context.Context, clientInfo *gRPC.Info) (*gRPC.Time, error) {
	//some code here:
	time := time.Now().GoString()
	return &gRPC.Time{Message: time}, nil
}
