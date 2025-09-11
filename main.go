package main

import (
	"log"
	"ygocdb-tui/internal/ui"
)

func main() {
	if err := ui.Start(); err != nil {
		log.Fatal(err)
	}
}