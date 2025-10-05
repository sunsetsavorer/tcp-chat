package server

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/sunsetsavorer/tcp-chat.git/config"
)

type ChatServer struct {
	messages chan string
	users    map[string]net.Conn
	mu       sync.RWMutex
	config   *config.Config
}

func New(config *config.Config) *ChatServer {

	return &ChatServer{
		messages: make(chan string, config.Server.MessagesBufferSize),
		users:    make(map[string]net.Conn),
		config:   config,
	}
}

func (s *ChatServer) Run() error {

	defer close(s.messages)

	listener, err := net.Listen("tcp", s.config.Server.Address)
	if err != nil {
		return err
	}

	fmt.Printf("listener was runned on %s\n", s.config.Server.Address)

	for range s.config.Server.ChatReceiversCount {
		go s.chatReceiver()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *ChatServer) handleConnection(conn net.Conn) {

	defer conn.Close()

	name := conn.RemoteAddr().String()
	isRegistered := false

	fmt.Printf("%s connected\n", name)
	s.sendToUser(
		conn,
		fmt.Sprintf("| Welcome to TCP chat %s, enter your nickname:", name),
	)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {

		text := scanner.Text()

		if !isRegistered {

			if text == "" {
				s.sendToUser(
					conn,
					"| Nickname can't be empty!!!",
				)
				continue
			}

			if _, ok := s.getConnectionByNickname(text); ok {
				s.sendToUser(
					conn,
					"| This nickname already taken!!!",
				)
				continue
			}

			isRegistered = true
			name = text

			s.sendToUser(
				conn,
				"| Use /help command to see list of available commands",
			)

			s.addUser(name, conn)
			s.sendToRoom(fmt.Sprintf("[ %s joined to chat ]", name))
			continue
		}

		if text == "Exit" {
			s.sendToUser(
				conn,
				fmt.Sprintf("| Bye %s!", name),
			)

			s.deleteUser(name)

			s.sendToRoom(fmt.Sprintf("[ %s left the chat ]", name))

			fmt.Printf("%s disconnected\n", name)

			break
		} else if text != "" {
			fmt.Printf("%s enters: %s\n", name, text)

			s.sendToRoom(fmt.Sprintf("> %s: %s", name, text))
		}
	}
}

func (s *ChatServer) sendToRoom(message string) {

	s.messages <- message
}

func (s *ChatServer) sendToUser(conn net.Conn, message string) {

	fmt.Fprintf(
		conn,
		"%s\r\n",
		message,
	)
}

func (s *ChatServer) chatReceiver() {

	for message := range s.messages {
		for _, conn := range s.users {
			fmt.Fprintf(conn, "%s\r\n", message)
		}
	}
}

func (s *ChatServer) getConnectionByNickname(nickname string) (net.Conn, bool) {

	s.mu.RLock()
	conn, ok := s.users[nickname]
	s.mu.RUnlock()

	return conn, ok
}

func (s *ChatServer) addUser(nickname string, conn net.Conn) {

	s.mu.Lock()
	s.users[nickname] = conn
	s.mu.Unlock()
}

func (s *ChatServer) deleteUser(nickname string) {

	s.mu.Lock()
	delete(s.users, nickname)
	s.mu.Unlock()
}
