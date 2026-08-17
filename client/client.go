package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var passname string
var ip string
var address string

func main() {

	for {
		fmt.Print("Input IP: ")
		_, err := fmt.Scan(&ip)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("\ngoodbye")
				return
			}
			fmt.Printf("scan error: %s\n", err)
			continue
		}

		address = fmt.Sprintf("%s:6661", ip)

		handleConnection()
	}
}

func handleConnection() {
	client, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("cannot connect: %s\n", err)
		return
	}
	defer client.Close()

	fmt.Printf("[+] Opened connection with %s\n", client.RemoteAddr())

	fmt.Print("Enter passname: ")
	_, err = fmt.Scan(&passname)
	if err != nil {
		fmt.Printf("scan error: %s\n", err)
		return
	}

	if _, err := client.Write([]byte(passname)); err != nil {
		fmt.Printf("cannot send passname: %s\n", err)
		return
	}

	if err := registerOrLogin(client); err != nil {
		fmt.Printf("registration failed: %s\n", err)
		return
	}

	done := make(chan struct{})
	go readLoop(client, done)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-done:
			return
		default:
		}

		message := scanner.Text()
		if message == "" {
			continue
		}

		if _, err := client.Write([]byte(message)); err != nil {
			fmt.Printf("cannot send message: %s\n", err)
			return
		}
	}
}

// registerOrLogin handles the server's response to the passname packet: a
// new passname gets "::REG::", prompting for a machine name; an existing
// one goes straight to "::OK::". Either way it blocks until "::OK::" arrives.
func registerOrLogin(client net.Conn) error {
	buf := make([]byte, 1024)

	n, err := client.Read(buf)
	if err != nil {
		return err
	}
	response := string(buf[:n])

	if response == "::REG::" {
		fmt.Print("Name your main machine: ")
		var machineName string
		if _, err := fmt.Scan(&machineName); err != nil {
			return err
		}

		if _, err := client.Write([]byte(machineName)); err != nil {
			return err
		}

		n, err = client.Read(buf)
		if err != nil {
			return err
		}
		response = string(buf[:n])
	}

	if response != "::OK::" {
		return fmt.Errorf("unexpected server response: %q", response)
	}

	fmt.Println("[+] Logged in successfully")
	return nil
}

// readLoop prints everything the server sends until the connection closes,
// then closes done so the input loop in handleConnection stops.
func readLoop(client net.Conn, done chan struct{}) {
	defer close(done)

	buf := make([]byte, 1024)
	for {
		n, err := client.Read(buf)
		if err != nil {
			fmt.Printf("\n[-] Connection closed: %s\n", err)
			return
		}

		message := strings.TrimSpace(string(buf[:n]))
		if message == "" {
			continue
		}
		fmt.Printf("< %s\n", message)
	}
}
