package outputfilegeneration

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"

	"github.com/pvelati/cluster-config-gen/types"
)

func GenerateFromTemplate(templateFile, outputFile string, cluster types.Cluster) {
	funcMap := template.FuncMap{
		"add1": func(i int) int { return i + 1 },
	}

	tmpl, err := template.New(filepath.Base(templateFile)).Funcs(funcMap).ParseFiles(templateFile)
	if err != nil {
		log.Fatalf("Errore nel parsare il template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cluster); err != nil {
		log.Fatalf("Errore nell'esecuzione del template: %v", err)
	}

	writeToFile(outputFile, buf.String())
}
