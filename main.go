package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var (
	messages chan string         = make(chan string, 5)
	users    map[string]net.Conn = make(map[string]net.Conn)
	mu       sync.RWMutex
)

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	fmt.Printf("listener was runned on :8080\n")

	for range 3 {
		go chatReceiver()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	name := conn.RemoteAddr().String()
	isRegistered := false

	fmt.Printf("%s connected\n", name)
	sendToUser(
		conn,
		fmt.Sprintf("Welcome to TCP chat %s, enter your nickname:", name),
	)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {

		text := scanner.Text()

		if !isRegistered {

			if text == "" {
				sendToUser(
					conn,
					"Nickname can't be empty!!!",
				)
				continue
			}

			if _, isUserExists := getConnectionByNickname(text); isUserExists {
				sendToUser(
					conn,
					"This nickname already taken!!!",
				)
				continue
			}

			isRegistered = true
			name = text

			addUser(name, conn)
			sendToRoom(fmt.Sprintf("%s joined to chat", name))
			continue
		}

		if text == "Exit" {
			sendToUser(
				conn,
				fmt.Sprintf("Bye %s!", name),
			)

			deleteUser(name)

			sendToRoom(fmt.Sprintf("%s left the chat", name))

			fmt.Printf("%s disconnected\n", name)

			break
		} else if text != "" {
			fmt.Printf("%s enters: %s\n", name, text)

			sendToRoom(fmt.Sprintf("%s: %s", name, text))
		}
	}
}

func sendToRoom(message string) {

	messages <- message
}

func sendToUser(conn net.Conn, message string) {

	fmt.Fprintf(
		conn,
		"%s\r\n",
		message,
	)
}

func chatReceiver() {

	for message := range messages {
		for _, conn := range users {
			fmt.Fprintf(conn, "%s\r\n", message)
		}
	}
}

func getConnectionByNickname(nickname string) (net.Conn, bool) {

	mu.RLock()
	conn, ok := users[nickname]
	mu.RUnlock()

	return conn, ok
}

func addUser(nickname string, conn net.Conn) {

	mu.Lock()
	users[nickname] = conn
	mu.Unlock()
}

func deleteUser(nickname string) {

	mu.Lock()
	delete(users, nickname)
	mu.Unlock()
}
