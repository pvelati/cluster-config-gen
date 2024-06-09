package outputfilegeneration

import (
	"fmt"
	"log"

	"github.com/pvelati/cluster-config-gen/types"
	"gopkg.in/yaml.v3"
)

func GenerateInventoryYAML(
	domain string, // FIXME: remove me
	outputFile string,
	config types.Config,
) {
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
