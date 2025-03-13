// client.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// Handle CTRL+C signal
	setupCleanExit()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Read messages from server
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("\nLost connection to server")
				os.Exit(0)
				return
			}
			// Clear the current input line, print the message, then reprint the prompt
			fmt.Print("\r\033[K") // Clear the current line
			fmt.Println(strings.TrimSpace(message))
			fmt.Print("You: ")
		}
	}()

	// Display initial prompt
	fmt.Print("You: ")

	// Send messages to server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()

		// If user types "quit", exit the program
		if message == "quit" {
			os.Exit(0)
		}

		// Send message to server
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}

		// Print a new prompt
		fmt.Print("You: ")
	}
}

// Handle CTRL+C signal to exit program immediately without confirmation
func setupCleanExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()
}
