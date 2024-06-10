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

func validateClusters(clusters []Cluster) error {
	nameSet := make(map[string]struct{})

	for _, cluster := range clusters {
		if cluster.Name == "" {
			return fmt.Errorf("name cannot be empty")
		}
		if _, exists := nameSet[cluster.Name]; exists {
			return fmt.Errorf("duplicate cluster name: %s", cluster.Name)
		}
		nameSet[cluster.Name] = struct{}{}

		if cluster.NumMaster != 1 && cluster.NumMaster != 3 {
			return fmt.Errorf("the number of masters in the %s cluster must be 1 or 3", cluster.Name)
		}
		if cluster.NumWorker < 1 {
			return fmt.Errorf("the %s cluster must have at least 1 worker", cluster.Name)
		}
		if cluster.MasterBaseVmid <= 0 {
			return fmt.Errorf("master_base_vmid should be a positive integer")
		}
		if cluster.MasterLastOctet <= 0 || cluster.MasterLastOctet >= 256 {
			return fmt.Errorf("master_last_octet should be between 1 and 255")
		}
		if cluster.MasterGateway <= 0 || cluster.MasterGateway >= 256 {
			return fmt.Errorf("master_gateway should be between 1 and 255")
		}

		// Worker fallback to Master values if empty
		if cluster.WorkerBaseVmid == 0 {
			cluster.WorkerBaseVmid = cluster.MasterBaseVmid
		}
		if cluster.WorkerAddressSansLastOctet == "" {
			cluster.WorkerAddressSansLastOctet = cluster.MasterAddressSansLastOctet
		}
		if cluster.WorkerLastOctet == 0 {
			cluster.WorkerLastOctet = cluster.MasterLastOctet
		}
		if cluster.WorkerGateway == 0 {
			cluster.WorkerGateway = cluster.MasterGateway
		}
		if cluster.WorkerDomain == "" {
			cluster.WorkerDomain = cluster.MasterDomain
		}

		if cluster.WorkerBaseVmid <= 0 {
			return fmt.Errorf("worker_base_vmid should be a positive integer")
		}
		if cluster.WorkerLastOctet <= 0 || cluster.WorkerLastOctet >= 256 {
			return fmt.Errorf("worker_last_octet should be between 1 and 255")
		}
		if cluster.WorkerGateway <= 0 || cluster.WorkerGateway >= 256 {
			return fmt.Errorf("worker_gateway should be between 1 and 255")
		}
	}
	return nil
}

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

		// Genera IP per i master
		for masterNodeIndex := 0; masterNodeIndex < cluster.NumMaster; masterNodeIndex++ {
			nodeNumber := masterNodeIndex + 1
			lastIpDigit := cluster.MasterLastOctet + masterNodeIndex
			host := fmt.Sprintf("k8s-%s-master-%d", cluster.Name, nodeNumber)
			ip := fmt.Sprintf("%s.%d", cluster.MasterAddressSansLastOctet, lastIpDigit)
			gateway := fmt.Sprintf("%s.%d", cluster.MasterAddressSansLastOctet, cluster.MasterGateway)
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
			gateway := fmt.Sprintf("%s.%d", cluster.WorkerAddressSansLastOctet, cluster.WorkerGateway)
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
