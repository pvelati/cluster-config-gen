package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Cluster rappresenta una configurazione di cluster Kubernetes.
type Cluster struct {
	Name            string   `yaml:"name"`
	NumMaster       int      `yaml:"num_master"`
	NumWorker       int      `yaml:"num_worker"`
	VIP             bool     `yaml:"vip"`
	MasterLastOctet int      `yaml:"master_last_octet"`
	WorkerLastOctet int      `yaml:"worker_last_octet"`
	MasterIPs       []string `yaml:"master_ips"`
	WorkerIPs       []string `yaml:"worker_ips"`
}

// Config rappresenta la struttura del file di configurazione.
type Config struct {
	Clusters []Cluster `yaml:"clusters"`
}

// Definire i primi 3 ottetti delle subnet
const masterSubnet = "192.168.0."
const workerSubnet = "192.168.1."
const domain = "home.lab"

func main() {
	// Leggi il file di configurazione YAML
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Errore nella lettura del file di configurazione: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Errore nel fare l'unmarshal dello YAML: %v", err)
	}

	// Genera gli indirizzi IP per ogni cluster
	for i, cluster := range config.Clusters {
		// Verifica che il numero di master sia 1 o 3
		if cluster.NumMaster != 1 && cluster.NumMaster != 3 {
			log.Fatalf("Il numero di master nel cluster %s deve essere 1 oppure 3.", cluster.Name)
		}

		// Verifica che ci sia almeno 1 worker
		if cluster.NumWorker < 1 {
			log.Fatalf("Il cluster %s deve avere almeno 1 worker.", cluster.Name)
		}

		// Genera IP per i master
		for j := 0; j < cluster.NumMaster; j++ {
			ip := fmt.Sprintf("%s%d", masterSubnet, cluster.MasterLastOctet+j)
			config.Clusters[i].MasterIPs = append(config.Clusters[i].MasterIPs, ip)
		}

		// Genera IP per i worker
		for j := 0; j < cluster.NumWorker; j++ {
			ip := fmt.Sprintf("%s%d", workerSubnet, cluster.WorkerLastOctet+j)
			config.Clusters[i].WorkerIPs = append(config.Clusters[i].WorkerIPs, ip)
		}
	}

	// Genera l'inventario YAML per Ansible
	generateInventoryYAML("ansible/inventory.yaml", config)

	// Genera il file di risorse Terraform per ogni cluster
	for _, cluster := range config.Clusters {
		generateFromTemplate("templates/terraform_template.tf.tmpl", fmt.Sprintf("terraform/%s_resources.tf", cluster.Name), cluster)
	}
}

func generateInventoryYAML(outputFile string, config Config) {
	inventory := make(map[string]interface{})

	for _, cluster := range config.Clusters {
		masterGroup := fmt.Sprintf("%s_master", cluster.Name)
		workerGroup := fmt.Sprintf("%s_worker", cluster.Name)

		inventory[masterGroup] = map[string]interface{}{
			"hosts": make(map[string]interface{}),
		}
		inventory[workerGroup] = map[string]interface{}{
			"hosts": make(map[string]interface{}),
		}

		for i, ip := range cluster.MasterIPs {
			host := fmt.Sprintf("k8s-%s-master-%d.%s", cluster.Name, i+1, domain)
			inventory[masterGroup].(map[string]interface{})["hosts"].(map[string]interface{})[host] = map[string]string{"ansible_host": ip}
		}
		for i, ip := range cluster.WorkerIPs {
			host := fmt.Sprintf("k8s-%s-worker-%d.%s", cluster.Name, i+1, domain)
			inventory[workerGroup].(map[string]interface{})["hosts"].(map[string]interface{})[host] = map[string]string{"ansible_host": ip}
		}

		clusterGroup := map[string]interface{}{
			"children": map[string]interface{}{
				masterGroup: nil,
				workerGroup: nil,
			},
		}

		inventory[cluster.Name] = clusterGroup
	}

	data, err := yaml.Marshal(inventory)
	if err != nil {
		log.Fatalf("Errore nel fare il marshal dello YAML: %v", err)
	}

	writeToFile(outputFile, string(data))
}

func generateFromTemplate(templateFile, outputFile string, cluster Cluster) {
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
