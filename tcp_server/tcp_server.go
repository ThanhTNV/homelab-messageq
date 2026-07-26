package tcp_server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThanhTNV/homelab-messageq/message_broker"
)

func Bootstrap() {
	// ======== Initialize a network listener

	// Create a TCP listener on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		// If there is an error creating the listener, panic and print the error
		// Because Golang does not have exceptions, we use panic to handle unexpected errors
		panic(err)
	}
	// Postpone closing the listener until the function returns
	defer listener.Close()

	fmt.Println("Server listening on port 8080...")

	// ======== Channel to handle OS signals for graceful shutdown

	// creates a buffered channel that can receive operating system signals
	sigChan := make(chan os.Signal, 1)
	// Notify the channel "sigChan" when an interrupt or termination signal is received
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start a goroutine (asynchronous function) to handle graceful shutdown when a signal is received
	go func() {
		// block until a signal arrives, and then handle it
		<-sigChan
		/* Close the listener to stop accepting new connections
		Because after os.Exit(0), the program will terminate immediately,
		and the deferred listener.Close() will not be executed, so we need to close it here explicitly.
		*/
		listener.Close()
		// Print must be call before os.Exit(0) because after os.Exit(0), the program will terminate immediately
		fmt.Println("\nShutting down the server...")
		// Exit the program with a status code of 0, indicating successful termination
		os.Exit(0)
	}()

	// ======== Accept incoming connections in a loop

	// var messageQueue = make(chan []byte, 100) // Buffered channel to hold messages (100 messages max)
	// subscribers := []net.Conn{}               // Slice to hold all subscriber connections

	commands := make(chan message_broker.BrokerCommand, 100) // buffered channel for broker commands
	go message_broker.BrokerLoop(commands)                   // Start the broker loop

	// Continuously accept incoming connections
	for {
		// Accept a new connection from the listener
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			break
		}

		// For each new connection, we create a new goroutine to handle it concurrently
		// Subscribers must be a pointer to the slice of connections, so that we can modify it in the HandleConnection function

		// go message_broker.HandleConnection(conn, messageQueue, &subscribers)
		go message_broker.HandleConnectionV2(conn, commands)
	}
}
