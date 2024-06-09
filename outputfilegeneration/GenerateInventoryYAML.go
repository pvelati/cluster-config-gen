package outputfilegeneration

import (
	"log"

	"github.com/pvelati/cluster-config-gen/types"
	"gopkg.in/yaml.v3"
)

// GenerateInventoryYAML genera un file di inventario YAML per Ansible.
func GenerateInventoryYAML(
	outputFile string, // Percorso del file di output
	internalData types.InternalData, // Dati interni dei cluster per la generazione dell'inventario
) {
	// Mappa per memorizzare l'inventario
	inventory := make(map[string]interface{})

	// Iterazione sui cluster per generare i gruppi di host
	for _, cluster := range internalData.Clusters {
		// Aggiunta dei gruppi master e worker all'inventario
		inventory[cluster.AnsibleMasterGroup] = map[string]interface{}{
			"hosts": make(map[string]interface{}),
		}
		inventory[cluster.AnsibleWorkerGroup] = map[string]interface{}{
			"hosts": make(map[string]interface{}),
		}

		// Aggiunta degli host master all'inventario
		for _, hostInfo := range cluster.Masters {
			inventory[cluster.AnsibleMasterGroup].(map[string]interface{})["hosts"].(map[string]interface{})[hostInfo.Host] = map[string]string{"ansible_host": hostInfo.IP}
		}
		// Aggiunta degli host worker all'inventario
		for _, hostInfo := range cluster.Workers {
			inventory[cluster.AnsibleWorkerGroup].(map[string]interface{})["hosts"].(map[string]interface{})[hostInfo.Host] = map[string]string{"ansible_host": hostInfo.IP}
		}

		// Creazione di un gruppo per il cluster
		clusterGroup := map[string]interface{}{
			"children": map[string]interface{}{
				cluster.AnsibleMasterGroup: nil,
				cluster.AnsibleWorkerGroup: nil,
			},
		}

		inventory[cluster.Name] = clusterGroup
	}

	// Conversione dell'inventario in YAML
	data, err := yaml.Marshal(inventory)
	if err != nil {
		log.Fatalf("Errore nel fare il marshal dello YAML: %v", err)
	}

	// Scrittura dei dati YAML nel file di output
	writeToFile(outputFile, string(data))
}
