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
const baseVmID = 4000

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

	var internalData types.InternalData

	// Genera gli indirizzi IP per ogni cluster
	for _, cluster := range config.Clusters {
		var oneCluster types.InternalDataCluster

		oneCluster.Name = cluster.Name
		oneCluster.AnsibleMasterGroup = fmt.Sprintf("%s_master", cluster.Name)
		oneCluster.AnsibleWorkerGroup = fmt.Sprintf("%s_worker", cluster.Name)

		// Verifica che il numero di master sia 1 o 3
		if cluster.NumMaster != 1 && cluster.NumMaster != 3 {
			log.Fatalf("Il numero di master nel cluster %s deve essere 1 oppure 3.", cluster.Name)
		}

		// Verifica che ci sia almeno 1 worker
		if cluster.NumWorker < 1 {
			log.Fatalf("Il cluster %s deve avere almeno 1 worker.", cluster.Name)
		}

		// Genera IP per i master
		for masterNodeIndex := 0; masterNodeIndex < cluster.NumMaster; masterNodeIndex++ {
			nodeNumber := masterNodeIndex + 1
			lastIpDigit := cluster.MasterLastOctet + masterNodeIndex
			host := fmt.Sprintf("k8s-%s-master-%d.%s", cluster.Name, nodeNumber, domain)
			ip := fmt.Sprintf("%s%d", masterSubnet, lastIpDigit)
			oneCluster.Masters = append(oneCluster.Masters, types.InternalDataMaster{
				IP:                    ip,
				Host:                  host,
				ProxmoxVMID:           baseVmID + lastIpDigit,
				TerraformResourceName: fmt.Sprintf("%s_master_%d", cluster.Name, nodeNumber),
				ProxmoxVmName:         fmt.Sprintf("%s-master-%d", cluster.Name, nodeNumber),
				ProxmoxVmDescription:  fmt.Sprintf("master node of kubernetes cluster %s", cluster.Name),
				ProxmoxVmTags: []string{
					fmt.Sprintf("k8s-cluster-%s-vm", cluster.Name),
					fmt.Sprintf("k8s-cluster-%s-master", cluster.Name),
				},
			})
		}

		// Genera IP per i worker
		for workerNodeIndex := 0; workerNodeIndex < cluster.NumWorker; workerNodeIndex++ {
			nodeNumber := workerNodeIndex + 1
			lastIpDigit := cluster.WorkerLastOctet + workerNodeIndex
			host := fmt.Sprintf("k8s-%s-worker-%d.%s", cluster.Name, nodeNumber, domain)
			ip := fmt.Sprintf("%s%d", workerSubnet, lastIpDigit)
			oneCluster.Workers = append(oneCluster.Workers, types.InternalDataWorker{
				IP:                    ip,
				Host:                  host,
				ProxmoxVMID:           baseVmID + lastIpDigit,
				TerraformResourceName: fmt.Sprintf("%s_worker_%d", cluster.Name, nodeNumber),
				ProxmoxVmName:         fmt.Sprintf("%s-worker-%d", cluster.Name, nodeNumber),
				ProxmoxVmDescription:  fmt.Sprintf("worker node of kubernetes cluster %s", cluster.Name),
				ProxmoxVmTags: []string{
					fmt.Sprintf("k8s-cluster-%s-vm", cluster.Name),
					fmt.Sprintf("k8s-cluster-%s-worker", cluster.Name),
				},
			})
		}

		internalData.Clusters = append(internalData.Clusters, oneCluster)
	}

	// Genera l'inventario YAML per Ansible
	outputfilegeneration.GenerateInventoryYAML("ansible/inventory.yaml", internalData)

	// Genera il file di risorse Terraform per ogni cluster
	for _, internalDataCluster := range internalData.Clusters {
		outputfilegeneration.GenerateFromTemplate("templates/terraform_template.tf.tmpl", fmt.Sprintf("terraform/%s_resources.tf", internalDataCluster.Name), internalDataCluster)
	}
}
