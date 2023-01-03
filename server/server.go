package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	dict "github.com/bemillant/dsExam/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	dict.UnimplementedDictionaryServer
	port       int // ownPort
	dictionary map[string]string
	servers    map[int32]dict.DictionaryClient
	isLeader   bool
}

/*
	- Assigns ownPort to user-input int + 5001. Default inputs are 0, 1 and 2.
	- Sets log.out.
	- Creates listener on ownPort.
*/

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5001

	setLog()
	flag.Parse()

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	server := &Server{
		port:       int(ownPort),
		dictionary: make(map[string]string),
		isLeader:   false,
	}

	if server.port == 5001 {
		server.servers = make(map[int32]dict.DictionaryClient)
		server.isLeader = true
		portTwo := int32(5001) + 1
		var conn *grpc.ClientConn
		insecure := insecure.NewCredentials()
		conn, err := grpc.Dial(fmt.Sprintf(":%v", portTwo), grpc.WithTransportCredentials(insecure), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		fmt.Printf("succesfully dialed to the other server")
		c := dict.NewDictionaryClient(conn)
		server.servers[portTwo] = c
	}

	grpcServer := grpc.NewServer()
	dict.RegisterDictionaryServer(grpcServer, server)

	for {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}
}

func (s *Server) Add(ctx context.Context, RequestBid *dict.RequestAdd) (*dict.Ack, error) {
	log.Printf("Attempting to add " + RequestBid.Value + " to the definition of the word: " + RequestBid.Key + ". Per instructed by " + RequestBid.Name)
	key := RequestBid.GetKey()
	val := RequestBid.GetValue()

	s.dictionary[key] = val

	if s.dictionary[RequestBid.Key] == RequestBid.Value {
		log.Printf("Add function was a success")
		return &dict.Ack{
			Message: "Succesfully added the definition '" + val + "' to the word '" + key + "'.",
			Success: true,
		}, nil
	} else {
		log.Printf("Add function failed")
		return &dict.Ack{
			Message: "Could not update the dictionary",
			Success: false,
		}, nil
	}
}

func (s *Server) Read(ctx context.Context, RequestBid *dict.ReadRequest) (*dict.ReadOutcome, error) {
	fmt.Println("Recieved Read-request to the word " + RequestBid.Key)
	log.Printf("Recieved Read-request to the word " + RequestBid.Key)

	key := RequestBid.GetKey()
	value, err := s.dictionary[key]
	if !err {
		log.Printf("The word did not exist in the dictionary and the client is returned: 'Non'")
		return &dict.ReadOutcome{
			Status: "The requested word does not exist in the dictionary as of yet.",
			Value:  "Non",
		}, nil
	}

	log.Printf("The client is returned the definition of their word: " + key + " - " + value)
	return &dict.ReadOutcome{
		Status: "Successfully requested definition",
		Value:  value,
	}, nil
}

// Sets log output to file in project dir
func setLog() *os.File {
	// Clears the log.txt file when a new server is started
	if err := os.Truncate("log.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	// This connects to the log file/changes the output of the log informaiton to the log.txt file.
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}
