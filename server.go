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
	"os/signal"
	"path/filepath"
	"syscall"
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

func main() {

	dataFilePath, err := getDataFilePath()
	if err != nil {
		log.Fatalf("couldn't resolve data file path: %s", err)
	}
	fmt.Println(dataFilePath)

	world, err := storage.LoadWorld(dataFilePath)
	if err != nil {
		log.Fatalf("couldn't load world: %s", err)
	}

	// making a listener
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("couldn't run a server: %s", err)
	}
	defer listener.Close()

	go waitForShutdown(dataFilePath, world, listener)

	fmt.Printf("[+] Server is running on port %s\n", PORT)

	// making a loop receiving connections and handling them
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Println(err)
			continue
		}
		fmt.Printf("[+] New connection: %s\n", conn.RemoteAddr())

		go handleConnection(conn, world)
	}
}

// waitForShutdown blocks until an interrupt/terminate signal arrives, saves
// the world to dataFilePath, then exits the process.
func waitForShutdown(dataFilePath string, world *objects.World, listener net.Listener) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\n[+] Shutting down, saving world...")
	if err := storage.SaveWorld(dataFilePath, world); err != nil {
		fmt.Printf("couldn't save world: %s\n", err)
	}

	listener.Close()
	os.Exit(0)
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
	var passname string
	var machineName string
	var playerMachine objects.Machine
	var player objects.Player

	// fetching passname as the first packet sent
	// initializing player
	n, err := conn.Read(buf)
	if err != nil {
		handleReceiveError(err)
		return
	}

	passname = string(buf[:n])
	fmt.Printf("Joined player (passname): %s\n", passname)

	if _, exists := world.Players[passname]; !exists {
		conn.Write([]byte("::REG::"))

		n, err = conn.Read(buf)
		if err != nil {
			fmt.Printf("cannot read machine name: %s", err)
			return
		}

		machineName = string(buf[:n])
		//fmt.Printf("machineName: %s\n", machineName)
		if err := world.NewPlayer(passname, machineName); err != nil {
			fmt.Printf("cannot register player: %s\n", err)
			return
		}

		player = world.Players[passname]
	} else {
		player = world.Players[passname]
		playerMachine = world.Machines[player.MainMachineID]
		machineName = playerMachine.MachineName
	}

	conn.Write([]byte("::OK::"))

	connHandler, err := objects.NewConnectionHandler(world, player.MainMachineID)
	var currMachine = world.Machines[connHandler.CurrentMachineID]

	for {
		n, err = conn.Read(buf)
		if err != nil {
			handleReceiveError(err)
			break
		}

		message := buf[:n]

		currMachine = world.Machines[connHandler.CurrentMachineID]
		fmt.Printf("%s@%s - %s\n", currMachine.PublicIp, machineName, string(message))
		conn.Write(message)
	}
}
