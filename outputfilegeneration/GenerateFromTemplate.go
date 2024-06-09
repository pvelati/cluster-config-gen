package outputfilegeneration

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"

	"github.com/pvelati/cluster-config-gen/types"
)

// GenerateFromTemplate genera un file utilizzando un template e dati interni specifici per un cluster.
func GenerateFromTemplate(
	templateFile string, // Percorso del file template
	outputFile string, // Percorso del file di output
	internalDataCluster types.InternalDataCluster, // Dati del cluster per popolare il template
) {
	// Parsing del template
	tmpl, err := template.New(filepath.Base(templateFile)).ParseFiles(templateFile)
	if err != nil {
		log.Fatalf("Errore nel parsare il template: %v", err)
	}

	var buf bytes.Buffer
	// Esecuzione del template e scrittura del risultato in un buffer
	if err := tmpl.Execute(&buf, internalDataCluster); err != nil {
		log.Fatalf("Errore nell'esecuzione del template: %v", err)
	}

	// Scrittura del contenuto del buffer nel file di output
	writeToFile(outputFile, buf.String())
}
