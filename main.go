package main

import "0x7266/go_server/internal"

func main() {
	server := internal.NewAPIServer(":3333")
	server.Run()
}
