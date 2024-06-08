// config/template.go
package config

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
)

// generateFromTemplate genera il file da un template.
func generateFromTemplate(templateFile, outputFile string, cluster Cluster) error {
	funcMap := template.FuncMap{
		"add1": func(i int) int { return i + 1 },
	}

	tmpl, err := template.New(filepath.Base(templateFile)).Funcs(funcMap).ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("Errore nel parsare il template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cluster); err != nil {
		return fmt.Errorf("Errore nell'esecuzione del template: %w", err)
	}

	if err := writeToFile(outputFile, buf.String()); err != nil {
		return fmt.Errorf("Errore nella scrittura del file: %w", err)
	}

	return nil
}
