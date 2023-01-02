package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	DB "github.com/bemillant/dsExam/grpc"
	"google.golang.org/grpc"
)

type Server struct {
	DB.UnimplementedAuctionServer
	port  int // ownPort
	dbMap map[int32]int32
}

/*
	- Assigns ownPort to user-input int + 5001. Default inputs are 0, 1 and 2.
	- Sets log.out.
	- Creates listener on ownPort.
	- Starts serving, and registers server.
	- Starts auction-timer (a bit stupidly, they should have begun at exactly the same time)
	- Announces winner to servers.

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
		port:  int(ownPort),
		dbMap: make(map[int32]int32),
	}

	grpcServer := grpc.NewServer()
	DB.RegisterAuctionServer(grpcServer, server)

	//go server.handleTime()

	for {

		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}

	}
}

func (s *Server) Put(ctx context.Context, RequestPut *DB.RequestPut) (*DB.Ack, error) {

	//log.Printf(RequestPut.Message)
	incomingKey := RequestPut.Key
	incomingValue := RequestPut.Value

	s.dbMap[incomingKey] = incomingValue

	if s.dbMap[incomingKey] == incomingValue {
		fmt.Println("added value to map")
		log.Print("The key/value pair: ")
		log.Print(incomingKey)
		log.Print(incomingValue)
		log.Print(" Has been added to the DB")
		return &DB.Ack{
			Message: "The Database has been updated.",
			Success: true,
		}, nil
	} else {
		return &DB.Ack{
			Message: "The Database could not be updated.",
			Success: false,
		}, nil
	}
}

func (s *Server) Get(ctx context.Context, getRequest *DB.ValueRequest) (*DB.Outcome, error) {

	isKeyInDB := false
	for k := range s.dbMap {
		if k == getRequest.Key {
			isKeyInDB = true
		}
	}

	if isKeyInDB {
		log.Print("A client is getting the following value returned")
		log.Print(s.dbMap[getRequest.GetKey()])
		return &DB.Outcome{
			Value: s.dbMap[getRequest.Key],
		}, nil
	} else {
		log.Print("A client is getting a value of a key that does not exist yet and is therefore returned the value '0'")
		return &DB.Outcome{
			Value: 0,
		}, nil
	}
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
