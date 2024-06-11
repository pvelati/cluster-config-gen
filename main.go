package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pvelati/cluster-config-gen/outputfilegeneration"
	"github.com/pvelati/cluster-config-gen/types"

	"gopkg.in/yaml.v3"
)

func main() {
	// Leggi il file di configurazione YAML
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config types.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Validate the clusters
	if err := validateClusters(config.Clusters); err != nil {
		log.Fatalf("Validation error: %v", err)
	}

	var internalData types.InternalData

	// Genera gli indirizzi IP per ogni cluster
	for _, cluster := range config.Clusters {
		var oneCluster types.InternalDataCluster

		oneCluster.Name = strings.ReplaceAll(cluster.Name, "-", "_")
		oneCluster.AnsibleMasterGroup = strings.ReplaceAll(fmt.Sprintf("%s_master", cluster.Name), "-", "_")
		oneCluster.AnsibleWorkerGroup = strings.ReplaceAll(fmt.Sprintf("%s_worker", cluster.Name), "-", "_")

		if cluster.MasterHa {
			cluster.NumMaster = 3
		} else {
			cluster.NumMaster = 1
		}

		// Genera IP per i master
		for masterNodeIndex := 0; masterNodeIndex < cluster.NumMaster; masterNodeIndex++ {
			nodeNumber := masterNodeIndex + 1
			lastIpDigit := cluster.MasterLastOctet + masterNodeIndex
			host := fmt.Sprintf("k8s-%s-master-%d", cluster.Name, nodeNumber)
			ip := fmt.Sprintf("%s.%d", cluster.MasterAddressSansLastOctet, lastIpDigit)
			gateway := fmt.Sprintf("%s.%d", cluster.MasterAddressSansLastOctet, cluster.MasterGatewayLastOctet)
			oneCluster.Masters = append(oneCluster.Masters, types.InternalDataMaster{
				IP:                    ip,
				Gateway:               gateway,
				Host:                  host,
				Domain:                cluster.MasterDomain,
				ProxmoxVMID:           cluster.MasterBaseVmid + lastIpDigit,
				TerraformResourceName: strings.ReplaceAll(host, "-", "_"),
				ProxmoxVmName:         strings.ReplaceAll(host, "_", "-"),
				ProxmoxVmDescription:  fmt.Sprintf("cluster %s - master node %d", cluster.Name, nodeNumber),
				ProxmoxVmTags: []string{
					strings.ReplaceAll(fmt.Sprintf("k8s_%s", cluster.Name), "-", "_"),
				},
			})
		}

		// Genera IP per i worker
		for workerNodeIndex := 0; workerNodeIndex < cluster.NumWorker; workerNodeIndex++ {
			nodeNumber := workerNodeIndex + 1
			lastIpDigit := cluster.WorkerLastOctet + workerNodeIndex
			host := fmt.Sprintf("k8s-%s-worker-%d", cluster.Name, nodeNumber)
			ip := fmt.Sprintf("%s.%d", cluster.WorkerAddressSansLastOctet, lastIpDigit)
			gateway := fmt.Sprintf("%s.%d", cluster.WorkerAddressSansLastOctet, cluster.WorkerGatewayLastOctet)
			oneCluster.Workers = append(oneCluster.Workers, types.InternalDataWorker{
				IP:                    ip,
				Gateway:               gateway,
				Host:                  host,
				Domain:                cluster.WorkerDomain,
				ProxmoxVMID:           cluster.WorkerBaseVmid + lastIpDigit,
				TerraformResourceName: strings.ReplaceAll(host, "-", "_"),
				ProxmoxVmName:         strings.ReplaceAll(host, "_", "-"),
				ProxmoxVmDescription:  fmt.Sprintf("cluster %s - worker node %d", cluster.Name, nodeNumber),
				ProxmoxVmTags: []string{
					strings.ReplaceAll(fmt.Sprintf("k8s_%s", cluster.Name), "-", "_"),
				},
			})
		}

		internalData.Clusters = append(internalData.Clusters, oneCluster)
	}

	// Genera l'inventario YAML per Ansible
	outputfilegeneration.GenerateInventoryYAML("ansible/inventory.yaml", internalData)

	// Genera il file di risorse Terraform per ogni cluster
	for _, internalDataCluster := range internalData.Clusters {
		outputfilegeneration.GenerateTerraformResource(fmt.Sprintf("terraform/%s_resources.tf", internalDataCluster.Name), internalDataCluster)
	}
}
