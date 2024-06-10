package outputfilegeneration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// writeToFile scrive il contenuto fornito in un file nel percorso specificato.
func writeToFile(filename, content string) {
	// Estrai la directory dal percorso del file
	dir := filepath.Dir(filename)
	// Crea la directory se non esiste
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Scrivi il contenuto nel file
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	// Stampa un messaggio di successo
	fmt.Printf("File %s generated successfully.\n", filename)
}
