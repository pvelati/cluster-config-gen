package outputfilegeneration

import (
	"log"

	"github.com/pvelati/cluster-config-gen/types"
	"gopkg.in/yaml.v3"
)

func GenerateInventoryYAML(
	outputFile string,
	internalData types.InternalData,
) {
	inventory := make(map[string]interface{})

	for _, cluster := range internalData.Clusters {
		inventory[cluster.AnsibleMasterGroup] = map[string]interface{}{
			"hosts": make(map[string]interface{}),
		}
		inventory[cluster.AnsibleWorkerGroup] = map[string]interface{}{
			"hosts": make(map[string]interface{}),
		}

		for _, hostInfo := range cluster.Masters {
			inventory[cluster.AnsibleMasterGroup].(map[string]interface{})["hosts"].(map[string]interface{})[hostInfo.Host] = map[string]string{"ansible_host": hostInfo.IP}
		}
		for _, hostInfo := range cluster.Workers {
			inventory[cluster.AnsibleWorkerGroup].(map[string]interface{})["hosts"].(map[string]interface{})[hostInfo.Host] = map[string]string{"ansible_host": hostInfo.IP}
		}

		clusterGroup := map[string]interface{}{
			"children": map[string]interface{}{
				cluster.AnsibleMasterGroup: nil,
				cluster.AnsibleWorkerGroup: nil,
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
