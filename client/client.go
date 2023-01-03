package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	dict "github.com/bemillant/dsExam/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	name    string
	server  dict.DictionaryClient
	servers map[int32]dict.DictionaryClient
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

	setLog()

	client := &Client{
		name:    tempName,
		servers: make(map[int32]dict.DictionaryClient),
		server:  nil,
	}

	go handleClient(client)

	for {

	}

}

func (client *Client) addToDictionary(word string, def string) {
	wordDef := &dict.RequestAdd{
		Name:  client.name,
		Key:   word,
		Value: def,
	}
	ack, err := client.server.Add(context.Background(), wordDef)
	if err != nil {
		fmt.Println("An error occured, possible on the servers side. The add operation was a %v. Trying to establish connection to new server...", ack.Success)
		log.Println("An error occured, possible on the servers side. Trying to establish connection to new server...")
		client.getServerConnection(5002)
		fmt.Println("Try again!")
	} else {
		if ack.Success {
			fmt.Println("The operation was a: %v!", ack.Success)
		}
		log.Println(client.name + " got the following result from the add request: " + ack.GetMessage())
		fmt.Println(ack.GetMessage())
	}
}

func (client *Client) requestRead(word string) {

	reqRead := &dict.ReadRequest{
		Key: word,
	}

	ack, err := client.server.Read(context.Background(), reqRead)
	if err != nil {
		fmt.Println("An error occured, possible on the servers side. Trying to establish connection to new server...")
		log.Println("An error occured, possible on the servers side. Trying to establish connection to new server...")
		client.getServerConnection(5002)
		fmt.Println("Try again!")
	} else {
		log.Printf(client.name + " got the following result from the Read request: " + ack.GetValue())
		fmt.Printf(ack.GetValue())
	}

}

// handles clientinput during runtime
func (client *Client) sendMessage() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text()

		if input == "Add" {
			client.handleAdd()
		} else if input == "Read" {
			client.handleRead()
		} else {
			fmt.Println("Try typing 'Add' to add a word and a definition to the dictionary or 'Read' to read a definition of a word")
		}
	}
}

func (client *Client) handleAdd() {
	fmt.Println("Enter the word followed by a '-' and then the words definition:")
	scanner := bufio.NewScanner(os.Stdin)
	slice := make([]string, 0)
	for scanner.Scan() {
		input := scanner.Text()
		if !strings.Contains(input, "-") {
			fmt.Println("The input does not contain '-' and therefore does not know the difference between the word and the definition")
			break
		}
		inputSlice := strings.Split(input, "-")
		for _, v := range inputSlice {
			slice = append(slice, strings.Trim(v, " "))
		}
		client.addToDictionary(slice[0], slice[1])
		break
	}
}
func (client *Client) handleRead() {
	fmt.Println("Enter the word you wish to see the definition of:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		client.requestRead(input)
		break
	}
}

func handleClient(client *Client) {
	client.getServerConnection(5001)
	go client.sendMessage()
	for {

	}
}

// hardcoded method that connects to two different servers to ensure active replications
func (client *Client) getServerConnection(portNum int32) {
	port := portNum
	var conn *grpc.ClientConn

	fmt.Printf("Trying to dial: %v\n", port)
	insecure := insecure.NewCredentials()
	conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithTransportCredentials(insecure), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	fmt.Printf("--- "+client.name+" succesfully dialed to %v\n", port)
	log.Printf("--- "+client.name+" succesfully dialed to %v\n", port)

	c := dict.NewDictionaryClient(conn)
	client.server = c

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
