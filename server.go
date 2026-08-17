package main

import (
	"errors"
	"fmt"
	"hackergame/objects"
	"hackergame/storage"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

var PORT = ":6661"
var DATAFILE_NAME = "world.json"

func getDataFilePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(exePath)
	dataFile := filepath.Join(dir, DATAFILE_NAME)

	return dataFile, nil
}

func initWorld() (*objects.World, error) {
	dataFilePath, err := getDataFilePath()
	if err != nil {
		return nil, err
	}

	world, err := storage.LoadWorld(dataFilePath)
	if err != nil {
		return nil, err
	}

	return world, nil
}

func main() {

	world, err := initWorld()
	if err != nil {
		log.Fatalf("couldn't initialize world: %s", err)
	}

	// making a listener
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("couldn't run a server: %s", err)
	}
	defer listener.Close()

	fmt.Printf("[+] Server is running on port %s\n", PORT)

	// making a loop receiving connections and handling them
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("[+] New connection: %s\n", conn.RemoteAddr())

		go handleConnection(conn, world)
	}
}

func handleReceiveError(err error) {
	if !errors.Is(err, io.EOF) {
		fmt.Printf("error: %s\n", err)
	}
}

func handleConnection(conn net.Conn, world *objects.World) {
	defer fmt.Printf("[+] Closed connection with %s\n", conn.RemoteAddr())
	defer conn.Close()
	defer conn.Write([]byte("exit:0"))

	buf := make([]byte, 1024)

	// fetching nickname as the first packet sent
	// initializing player
	n, err := conn.Read(buf)
	if err != nil {
		handleReceiveError(err)
		return
	}
	nickname := string(buf[:n])
	fmt.Printf("Player's nickname: %s\n", nickname)

	if _, exists := world.Players[nickname]; !exists {
		conn.Write([]byte("::REG::"))

		n, err = conn.Read(buf)
		if err != nil {
			fmt.Printf("cannot read machine name: %s", err)
			return
		}

		machineName := string(buf[:n])
		if err := world.NewPlayer(nickname, machineName); err != nil {
			fmt.Printf("cannot register player: %s\n", err)
			return
		}
	}

	conn.Write([]byte("::OK::"))

	for {
		n, err = conn.Read(buf)
		if err != nil {
			handleReceiveError(err)
			break
		}

		packet := buf[:n]
		fmt.Printf("%s - %s\n", conn.RemoteAddr(), string(packet))
		conn.Write(packet)
	}
}
