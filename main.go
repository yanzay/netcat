package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	// Server
	listen     = flag.Bool("l", false, "Listen")
	host       = flag.String("h", "localhost", "Host")
	port       = flag.Int("p", 0, "Port")
	command    = flag.Bool("c", false, "Command server")
	fileServer = flag.Bool("f", false, "Server for file upload")
	// Client
	execute = flag.String("e", "", "Execute command")
	upload  = flag.String("u", "", "Upload file")
)

func main() {
	flag.Parse()
	if *listen {
		addr := fmt.Sprintf("%s:%d", *host, *port)
		err := startServer(addr, *command, *fileServer)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if len(flag.Args()) < 2 {
			fmt.Println("Hostname and port required")
			os.Exit(1)
		}
		serverHost := flag.Arg(0)
		serverPort := flag.Arg(1)
		addr := fmt.Sprintf("%s:%s", serverHost, serverPort)
		err := startClient(addr, *execute, *upload)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
