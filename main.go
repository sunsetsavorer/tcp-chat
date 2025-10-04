package main

import (
	"bufio"
	"fmt"
	"net"
)

/*

	TODO

	Создание пользователя
	Общение
	Выход из чата

*/

func main() {

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	fmt.Printf("listener was runned on :8080\n")

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

	fmt.Printf("%s connected\n", name)
	fmt.Fprintf(
		conn,
		"Welcome to TCP chat %s\n\r",
		name,
	)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {

		text := scanner.Text()

		if text == "Exit" {
			fmt.Printf("%s disconnected\n", name)
			fmt.Fprintf(
				conn,
				"Bye!\n\r",
			)

			break
		} else if text != "" {
			fmt.Printf("%s enters: %s\n", name, text)
			fmt.Fprintf(
				conn,
				"You enter: %s\n\r",
				text,
			)
		}
	}
}
