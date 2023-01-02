package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	DB "github.com/bemillant/dsExam/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	name    string
	servers map[int32]DB.AuctionClient
}

func main() {

	fmt.Println("Please enter a name")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanBytes)

	// Set client name:
	tempName := ""
	for scanner.Scan() {
		if scanner.Text() == "\n" {
			break
		} else {
			tempName += scanner.Text()
		}
	}

	// flag.Parse()
	setLog()

	client := &Client{
		name:    tempName,
		servers: make(map[int32]DB.AuctionClient),
	}

	go handleClient(client)

	for {

	}

}

// function to put value into DB
func (client *Client) Put(key int32, value int32) {
	put := &DB.RequestPut{
		Name:  client.name,
		Key:   key,
		Value: value,
	}

	// List of responses from replicas when asking for auction result:
	resultList := make([]string, 0, 3)

	for port, server := range client.servers {
		ack, err := server.Put(context.Background(), put)
		if err != nil {
			delete(client.servers, port)
			log.Printf(client.name + "lost connection to a server, operating number of servers are now " + strconv.Itoa(int(len(resultList))))
			// fmt.Printf("something went wrong in bid method: %v", err)
		} else {
			resultList = append(resultList, ack.GetMessage())
		}

	}
	fmt.Printf(resultList[0] + "\n")
}

// get function
func (client *Client) Get(key int32) {
	reqVal := &DB.ValueRequest{
		Key: key,
	}

	resultList := make([]int32, 0, 3)

	for port, server := range client.servers {
		outcome, err := server.Get(context.Background(), reqVal)
		if err != nil {
			delete(client.servers, port)
			log.Printf(client.name + "lost connection to a server, operating number of servers are now " + strconv.Itoa(int(len(resultList))))
			// fmt.Printf("something went wrong in requestHB method: %v", err)
		} else {
			resultList = append(resultList, outcome.GetValue())
		}

	}
	fmt.Println(resultList[0])
}

// handles clientinput
func (client *Client) sendMessage() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()

		if input == "Put" {
			client.handlePut()
		} else if input == "Get" {
			client.handleGet()
		} else {
			fmt.Println("Try 'Put' to put values into the DB or 'Get' to search for a value")
		}

		// //Checks if the client input is of type string og integer (Tries to convert/parse, if an error occurs = string is not parseable)
		// amount, err := strconv.ParseInt(input, 10, 32)
		// if err != nil {
		// 	//Tell client that the current bid is at ... and to make a bid, type an integer
		// 	client.requestHighestBid()
		// } else {
		// 	client.makeBid(int32(amount))
		// }
	}
}

func (client *Client) handlePut() {
	fmt.Println("Enter the key value pair you want to add with a space inbetween")
	scanner := bufio.NewScanner(os.Stdin)
	slice := make([]int32, 0)
	for scanner.Scan() {
		input := scanner.Text()
		inputSlice := strings.Split(input, " ")
		for _, v := range inputSlice {
			number, _ := strconv.Atoi(v)
			slice = append(slice, int32(number))
		}
		break
	}
	client.Put(slice[0], slice[1])
}
func (client *Client) handleGet() {
	fmt.Println("Enter the key of the value you want to recieve")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		key, _ := strconv.Atoi(input)

		client.Get(int32(key))

		break
	}
}

func handleClient(client *Client) {
	client.getServerConnection()

	go client.sendMessage()

	for {

	}

}

// hardcoded method that connects to three different servers to ensure active replications
func (client *Client) getServerConnection() {

	//Connection to 3 servers
	for i := 0; i < 3; i++ {

		port := int32(5001) + int32(i)
		var conn *grpc.ClientConn

		fmt.Printf("Trying to dial: %v\n", port)
		insecure := insecure.NewCredentials()
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithTransportCredentials(insecure), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}

		fmt.Printf("--- "+client.name+" succesfully dialed to %v\n", port)
		log.Printf("--- "+client.name+" succesfully dialed to %v\n", port)

		// defer conn.Close()
		c := DB.NewAuctionClient(conn)
		client.servers[port] = c
	}

	//If only one server use this:
	/*
		port := int32(5001)
		var conn *grpc.ClientConn

		fmt.Printf("Trying to dial: %v\n", port)
		insecure := insecure.NewCredentials()
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithTransportCredentials(insecure), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		fmt.Printf("--- "+client.name+" succesfully dialed to %v\n", port)
		log.Printf("--- "+client.name+" succesfully dialed to %v\n", port)

		// defer conn.Close()
		c := auction.NewAuctionClient(conn)
		client.servers[port] = c
	*/

}

// Sets log output to file in project dir
func setLog() *os.File {
	// This connects to the log file/changes the output of the log informaiton to the log.txt file.
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}
