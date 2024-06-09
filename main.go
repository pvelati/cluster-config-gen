package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pvelati/cluster-config-gen/outputfilegeneration"
	"github.com/pvelati/cluster-config-gen/types"

	"gopkg.in/yaml.v3"
)

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

	var config types.Config
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
	outputfilegeneration.GenerateInventoryYAML(domain, "ansible/inventory.yaml", config)

	// Genera il file di risorse Terraform per ogni cluster
	for _, cluster := range config.Clusters {
		outputfilegeneration.GenerateFromTemplate("templates/terraform_template.tf.tmpl", fmt.Sprintf("terraform/%s_resources.tf", cluster.Name), cluster)
	}
}
