package main

import (
	"math/rand"
	"time"

	"github.com/zlahrouni/loganizer/cmd"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	cmd.Execute()
}
