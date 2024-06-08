// main.go
package main

import (
	"log"

	"github.com/pvelati/pvelati/cluster-config-gen/config"
)

func main() {
	if err := config.Run(); err != nil {
		log.Fatalf("Errore durante l'esecuzione del programma: %v", err)
	}
}
