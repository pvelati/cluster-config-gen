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

		if cluster.Controlplane.Cluster {
			cluster.Controlplane.Num = 3
			oneCluster.Ha = true
			if cluster.Vip.Controlplane {
				oneCluster.HaIp = fmt.Sprintf("%s.%d", cluster.Controlplane.AddressSansLastOctet, cluster.Controlplane.LastOctet)
				oneCluster.HaFqdn = fmt.Sprintf("%s.%s", cluster.Controlplane.Domain, cluster.Controlplane.Domain)
				oneCluster.Vip = true
			} else {
				oneCluster.HaIp = cluster.Vip.ControlplaneIp
				oneCluster.HaFqdn = cluster.Vip.ControlplaneFqdn
				oneCluster.Vip = false
			}
		} else {
			cluster.Controlplane.Num = 1
			oneCluster.Ha = false
		}

		// Genera IP per i master
		for masterNodeIndex := 0; masterNodeIndex < cluster.Controlplane.Num; masterNodeIndex++ {
			nodeNumber := masterNodeIndex + 1
			lastIpDigit := cluster.Controlplane.LastOctet + masterNodeIndex
			if cluster.Controlplane.Cluster && cluster.Vip.Controlplane {
				lastIpDigit++
			}
			host := fmt.Sprintf("k8s-%s-master-%d", cluster.Name, nodeNumber)
			ip := fmt.Sprintf("%s.%d", cluster.Controlplane.AddressSansLastOctet, lastIpDigit)
			gateway := fmt.Sprintf("%s.%d", cluster.Controlplane.AddressSansLastOctet, cluster.Controlplane.GatewayLastOctet)
			oneCluster.Masters = append(oneCluster.Masters, types.InternalDataNode{
				IP:                    ip,
				Gateway:               gateway,
				Host:                  host,
				Domain:                cluster.Controlplane.Domain,
				Core:                  cluster.Controlplane.Core,
				Memory:                cluster.Controlplane.Memory,
				TerraformResourceName: strings.ReplaceAll(host, "-", "_"),
				ProxmoxVMID:           cluster.Controlplane.BaseVmid + lastIpDigit,
				ProxmoxVmName:         strings.ReplaceAll(host, "_", "-"),
				ProxmoxVmDescription:  fmt.Sprintf("cluster %s - master node %d", cluster.Name, nodeNumber),
				ProxmoxVmTags: []string{
					strings.ReplaceAll(fmt.Sprintf("k8s_%s", cluster.Name), "-", "_"),
				},
			})
		}

		// Genera IP per i worker
		for workerNodeIndex := 0; workerNodeIndex < cluster.Compute.Num; workerNodeIndex++ {
			nodeNumber := workerNodeIndex + 1
			lastIpDigit := cluster.Compute.LastOctet + workerNodeIndex
			host := fmt.Sprintf("k8s-%s-worker-%d", cluster.Name, nodeNumber)
			ip := fmt.Sprintf("%s.%d", cluster.Compute.AddressSansLastOctet, lastIpDigit)
			gateway := fmt.Sprintf("%s.%d", cluster.Compute.AddressSansLastOctet, cluster.Compute.GatewayLastOctet)
			oneCluster.Workers = append(oneCluster.Workers, types.InternalDataNode{
				IP:                    ip,
				Gateway:               gateway,
				Host:                  host,
				Domain:                cluster.Compute.Domain,
				Core:                  cluster.Compute.Core,
				Memory:                cluster.Compute.Memory,
				TerraformResourceName: strings.ReplaceAll(host, "-", "_"),
				ProxmoxVMID:           cluster.Compute.BaseVmid + lastIpDigit,
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

	// Genera i group_vars Ansible e risorse Terraform per ogni cluster
	for _, internalDataCluster := range internalData.Clusters {
		outputfilegeneration.GenerateGroupVarsYAML(fmt.Sprintf("ansible/group_vars/%s.yaml", internalDataCluster.Name), internalDataCluster)
		outputfilegeneration.GenerateTerraformResource(fmt.Sprintf("terraform/%s_resources.tf", internalDataCluster.Name), internalDataCluster)
	}
}
