package main

import (
	"0x7266/go_server/internal"
	"log"
)

func main() {
	server := internal.NewAPIServer(":3333")
	log.Fatalln(server.Run())
}
