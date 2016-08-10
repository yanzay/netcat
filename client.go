package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func startClient(addr string, execute string, upload string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("Can't connect to server: %s\n", err)
	}
	if len(execute) > 0 && len(upload) > 0 {
		return fmt.Errorf("Can't execute command and upload file simultaneously")
	}
	switch {
	case len(execute) > 0:
		err = commandClient(execute, conn)
	case len(upload) > 0:
		err = fileClient(upload, conn)
	default:
		err = defaultClient(conn)
	}
	return err
}

func defaultClient(conn io.ReadWriteCloser) error {
	go io.Copy(os.Stdout, conn)
	_, err := io.Copy(conn, os.Stdin)
	if err != nil {
		return err
	}
	return conn.Close()
}

func commandClient(command string, conn io.ReadWriteCloser) error {
	_, err := io.WriteString(conn, fmt.Sprintf("%s\n", command))
	if err != nil {
		return err
	}
	return defaultClient(conn)
}

func fileClient(filename string, conn io.WriteCloser) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	return uploadFile(stat.Name(), conn, file)
}

func uploadFile(name string, conn io.WriteCloser, file io.ReadCloser) error {
	_, err := io.WriteString(conn, fmt.Sprintf("%s\n", name))
	if err != nil {
		return err
	}
	_, err = io.Copy(conn, file)
	if err != nil {
		return err
	}
	err = conn.Close()
	if err != nil {
		return err
	}
	return file.Close()
}
