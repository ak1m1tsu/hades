package main

import (
	"log"

	"github.com/ak1m1tsu/hades/internal/app/hades"
)

func main() {
	if err := hades.New().WithLogger().WithAPI().Run(); err != nil {
		log.Fatal(err)
	}
}
