// config/config.go
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
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

// Run esegue il programma.
func Run() error {
	// Leggi il file di configurazione YAML
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("Errore nella lettura del file di configurazione: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("Errore nel fare l'unmarshal dello YAML: %w", err)
	}

	// Genera gli indirizzi IP per ogni cluster
	for i, cluster := range config.Clusters {
		// Verifica che il numero di master sia 1 o 3
		if cluster.NumMaster != 1 && cluster.NumMaster != 3 {
			return fmt.Errorf("Il numero di master nel cluster %s deve essere 1 oppure 3.", cluster.Name)
		}

		// Verifica che ci sia almeno 1 worker
		if cluster.NumWorker < 1 {
			return fmt.Errorf("Il cluster %s deve avere almeno 1 worker.", cluster.Name)
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
	if err := generateInventoryYAML("outputs/inventory.yaml", config); err != nil {
		return fmt.Errorf("Errore durante la generazione dell'inventario YAML: %w", err)
	}

	// Genera il file di risorse Terraform per ogni cluster
	for _, cluster := range config.Clusters {
		if err := generateFromTemplate("templates/terraform_template.tf.tmpl", fmt.Sprintf("outputs/%s_resources.tf", cluster.Name), cluster); err != nil {
			return fmt.Errorf("Errore durante la generazione del file di risorse Terraform: %w", err)
		}
	}

	return nil
}
