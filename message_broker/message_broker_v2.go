package message_broker

import (
	"fmt"
	"net"
)

// Commands for the broker loop
type BrokerCommand struct {
	add     net.Conn // Connection to add as a subscriber
	remove  net.Conn // Connection to remove from subscribers
	publish []byte   // Message to publish to all subscribers
}

// Broker loop: owns the subscriber slice
// This function runs in a separate goroutine and handles all commands sent to the broker.
func BrokerLoop(commands <-chan BrokerCommand) {
	subscribers := []net.Conn{}

	// Continuously listen for commands from the channel
	for cmd := range commands {
		// Handle the command based on its type (add, remove, publish)

		// Add a new subscriber
		if cmd.add != nil {
			subscribers = append(subscribers, cmd.add)
			fmt.Println("New subscriber added:", cmd.add.RemoteAddr())
		}

		// Remove a subscriber
		if cmd.remove != nil {
			for i, sub := range subscribers {
				if sub == cmd.remove {
					subscribers = append(subscribers[:i], subscribers[i+1:]...)
					fmt.Println("Subscriber removed:", sub.RemoteAddr())
					break
				}
			}
		}

		// Publish a message to all subscribers
		if cmd.publish != nil {
			for _, sub := range subscribers {
				_, err := sub.Write(cmd.publish)
				if err != nil {
					fmt.Println("Error writing to subscriber:", err)
				}
			}
		}
	}
}

// Handle a new connection
func HandleConnectionV2(conn net.Conn, commands chan<- BrokerCommand) {
	defer conn.Close()

	// Tell broker to add this subscriber by sending an add command to the commands channel
	commands <- BrokerCommand{add: conn}

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			// Tell broker to remove this subscriber by sending a remove command to the commands channel
			commands <- BrokerCommand{remove: conn}
			return
		}
		// Copy data and send publish command
		data := make([]byte, n)
		copy(data, buffer[:n])

		// Send the publish command to the broker loop
		commands <- BrokerCommand{publish: data}
	}
}
