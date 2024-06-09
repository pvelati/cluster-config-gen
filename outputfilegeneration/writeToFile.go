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
		log.Fatalf("Errore nella creazione della directory: %v", err)
	}

	// Scrivi il contenuto nel file
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		log.Fatalf("Errore nella scrittura del file: %v", err)
	}

	// Stampa un messaggio di successo
	fmt.Printf("File %s generato con successo.\n", filename)
}
