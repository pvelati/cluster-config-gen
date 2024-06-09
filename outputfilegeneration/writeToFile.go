package outputfilegeneration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func writeToFile(filename, content string) {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Errore nella creazione della directory: %v", err)
	}

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		log.Fatalf("Errore nella scrittura del file: %v", err)
	}

	fmt.Printf("File %s generato con successo.\n", filename)
}
