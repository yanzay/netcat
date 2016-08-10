package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

func startServer(addr string, command bool, fileServer bool) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Printf("Listening for connections on %s", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
		} else {
			go processClient(conn, command, fileServer)
		}
	}
}

func processClient(conn io.ReadWriteCloser, command bool, fileServer bool) error {
	if command && fileServer {
		return fmt.Errorf("Can't launch server in command and file mode simultaneously")
	}
	var err error
	switch {
	case command:
		err = commandProcessor(conn)
	case fileServer:
		err = fileProcessor(conn)
	default:
		err = defaultProcessor(conn)
	}
	return err
}

func defaultProcessor(conn io.ReadCloser) error {
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		return err
	}
	return conn.Close()
}

func commandProcessor(conn io.ReadWriteCloser) error {
	defer conn.Close()
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
	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}

func fileProcessor(conn io.ReadCloser) error {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	line = strings.TrimSpace(line)
	file, err := os.Create(line)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, conn)
	return err
}
