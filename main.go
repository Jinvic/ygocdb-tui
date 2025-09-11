package main

import (
	"log"
	"ygocdb-tui/ui"
)

func main() {
	if err := ui.Start(); err != nil {
		log.Fatal(err)
	}
}