package main

import (
	"agent/workers"
)

func main() {
	workers.Set()
	StartServer(8080)
}