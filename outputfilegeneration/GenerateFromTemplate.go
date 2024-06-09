package outputfilegeneration

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"

	"github.com/pvelati/cluster-config-gen/types"
)

func GenerateFromTemplate(
	templateFile string,
	outputFile string,
	internalDataCluster types.InternalDataCluster,
) {
	tmpl, err := template.New(filepath.Base(templateFile)).ParseFiles(templateFile)
	if err != nil {
		log.Fatalf("Errore nel parsare il template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, internalDataCluster); err != nil {
		log.Fatalf("Errore nell'esecuzione del template: %v", err)
	}

	writeToFile(outputFile, buf.String())
}
