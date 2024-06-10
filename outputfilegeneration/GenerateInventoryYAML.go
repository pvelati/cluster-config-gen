package outputfilegeneration

import (
	"log"

	"github.com/pvelati/cluster-config-gen/types"
	"gopkg.in/yaml.v3"
)

type ansibleInventoryGroupType struct {
	Children map[string]any                      `yaml:"children,omitempty"`
	Hosts    map[string]ansibleInventoryHostType `yaml:"hosts,omitempty"`
}

type ansibleInventoryHostType struct {
	AnsibleHost string
}

// GenerateInventoryYAML genera un file di inventario YAML per Ansible.
func GenerateInventoryYAML(
	outputFile string, // Percorso del file di output
	internalData types.InternalData, // Dati interni dei cluster per la generazione dell'inventario
) {
	// Mappa per memorizzare l'inventario
	inventory := map[string]ansibleInventoryGroupType{}

	// Iterazione sui cluster per generare i gruppi di host
	for _, cluster := range internalData.Clusters {
		// Aggiunta degli host master all'inventario
		inventoryAnsibleMasterGroupHosts := map[string]ansibleInventoryHostType{}
		for _, hostInfo := range cluster.Masters {
			inventoryAnsibleMasterGroupHosts[hostInfo.Host+"."+hostInfo.Domain] = ansibleInventoryHostType{
				AnsibleHost: hostInfo.IP,
			}
		}
		inventory[cluster.AnsibleMasterGroup] = ansibleInventoryGroupType{
			Hosts: inventoryAnsibleMasterGroupHosts,
		}

		// Aggiunta degli host worker all'inventario
		inventoryAnsibleWorkerGroupHosts := map[string]ansibleInventoryHostType{}
		for _, hostInfo := range cluster.Workers {
			inventoryAnsibleWorkerGroupHosts[hostInfo.Host+"."+hostInfo.Domain] = ansibleInventoryHostType{
				AnsibleHost: hostInfo.IP,
			}
		}
		inventory[cluster.AnsibleWorkerGroup] = ansibleInventoryGroupType{
			Hosts: inventoryAnsibleWorkerGroupHosts,
		}

		// Creazione di un gruppo per il cluster
		inventory[cluster.Name] = ansibleInventoryGroupType{
			Children: map[string]any{
				cluster.AnsibleMasterGroup: nil,
				cluster.AnsibleWorkerGroup: nil,
			},
		}
	}

	// Conversione dell'inventario in YAML
	data, err := yaml.Marshal(inventory)
	if err != nil {
		log.Fatalf("Errore nel fare il marshal dello YAML: %v", err)
	}

	// Scrittura dei dati YAML nel file di output
	writeToFile(outputFile, string(data))
}
