package main

import "workshop/internal/cmd"

//go:generate go tool sqlc generate

func main() {
	cmd.Run()
}
