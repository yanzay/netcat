package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

var (
	listen  = flag.Bool("l", false, "Listen")
	host    = flag.String("h", "localhost", "Host")
	port    = flag.Int("p", 0, "Port")
	command = flag.Bool("c", false, "Command server")
	execute = flag.String("e", "", "Execute command")
)

func main() {
	flag.Parse()
	if *listen {
		startServer()
		return
	}
	if len(flag.Args()) < 2 {
		fmt.Println("Hostname and port required")
		return
	}
	serverHost := flag.Arg(0)
	serverPort := flag.Arg(1)
	startClient(fmt.Sprintf("%s:%s", serverHost, serverPort))
}

func startServer() {
	addr := fmt.Sprintf("%s:%d", *host, *port)
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening for connections on %s", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
		} else {
			go processClient(conn)
		}
	}
}

func processClient(conn net.Conn) {
	if *command {
		err := launchCommand(conn)
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
	}
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		log.Println(err)
	}
	conn.Close()
}

func launchCommand(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	line = strings.TrimSpace(line)
	log.Printf("Command: %s\n", line)
	cmd := exec.Command(line)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(stdin, conn)
	go io.Copy(conn, stdout)
	go io.Copy(conn, stderr)
	return cmd.Run()
}

func startClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Can't connect to server: %s\n", err)
		return
	}
	if len(*execute) > 0 {
		cmd := fmt.Sprintf("%s\n", *execute)
		conn.Write([]byte(cmd))
	}
	go io.Copy(os.Stdout, conn)
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		fmt.Printf("Connection error: %s\n", err)
	}
}
