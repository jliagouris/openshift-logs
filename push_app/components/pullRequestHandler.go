package main

import (
	"fmt"
	"io"
	"net"
	"os"

	yaml "gopkg.in/yaml.v3"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "1296"
	CONN_TYPE = "tcp"
)

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error Listening", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

type T struct {
	promql_req string `yaml:"a,omitempty"`
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading:", err)
		}
	}

	fmt.Println(string(buf))

	t := T{}

	err = yaml.Unmarshal(buf, &t)
	fmt.Println(t)
	if err != nil {
		fmt.Println("could not parse yaml file")
	}
	// fmt.Printf("--- d dump:\n%s\n\n", string(d))

	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))

	// Close the connection when you're done with it.
	conn.Close()
}
