package main

import (
	"agent/starter"
	"agent/workers"
)

func main() {
	workers.Set()
	starter.StartAgent()
}