package message_broker

import (
	"fmt"
	"net"
	"sync"
)

var mu sync.Mutex

func HandleConnection(conn net.Conn, messageQueue chan []byte, pSubscribers *[]net.Conn) {
	/* Always postpone closing the connection until the function returns,
	to ensure resources are released properly
	*/
	defer conn.Close()

	// Subscriber must be able to receive messages from the message queue
	// as well as send messages to the message queue

	//======== Start a goroutine to read data from the connection and send it to the message queue
	go subscribeNewConsumer(conn, messageQueue, pSubscribers)

	//======== Start a goroutine to read data from the message queue and send it to the connection
	go publishMessageToSubscribers(messageQueue, *pSubscribers)
}

func subscribeNewConsumer(conn net.Conn, messageQueue chan []byte, pSubscribers *[]net.Conn) {
	//======== Add the new consumer to the list of subscribers

	// // Use a mutex to ensure that only one goroutine can modify the subscribers slice at a time
	mu.Lock()
	*pSubscribers = append(*pSubscribers, conn)
	// Unlock the mutex to allow other goroutines to access the subscribers slice
	mu.Unlock()

	buffer := make([]byte, 1024) // 1KB buffer
	for {
		// Read data from the connection into the buffer
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			break
		}
		// Copy the data read into a new slice to avoid overwriting the buffer in the next iteration
		data := make([]byte, n)
		copy(data, buffer[:n])
		messageQueue <- data
	}
}

func publishMessageToSubscribers(messageQueue chan []byte, subscribers []net.Conn) {
	for {
		// Read data from the message queue
		data := <-messageQueue
		// Send the data to all subscribers
		mu.Lock()
		for _, subscriber := range subscribers {
			_, err := subscriber.Write(data)
			if err != nil {
				fmt.Println("Error writing to subscriber:", err)
			}
		}
		mu.Unlock()
	}
}
