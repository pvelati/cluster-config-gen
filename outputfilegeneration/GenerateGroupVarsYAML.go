package outputfilegeneration

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/pvelati/cluster-config-gen/types"

	"gopkg.in/yaml.v3"
)

func GenerateGroupVarsYAML(
	outputFile string, // Percorso del file di output
	internalDataCluster types.InternalDataCluster, // Dati interni dei cluster per la generazione dell'inventario
) {

	// Creazione della mappa per memorizzare i dati YAML prendendoli in input
	yamlGroupVars := map[string]interface{}{
		"registry_mirror":            false,
		"container_registry_address": "192.168.10.100",
		"cluster_token":              fmt.Sprintf("%x", sha256.Sum256([]byte(internalDataCluster.Name)))[:15],
		"services_vip_address":       "192.168.10.120",
		"services_range_start":       "192.168.10.121",
		"services_range_end":         "192.168.10.122",
		"active_kubeconfig_file":     "kubeconfig_test.yaml",
		"github_token":               "aaaaaaaaaaaaaaaaaaa",
		"github_user":                "pvelati",
		"github_repository":          "flux-test",
		"github_repo_path":           "clusters/k01",
	}

	// Aggiunta di cluster_vip_fqdn e cluster_vip_address solo se internalDataCluster.Ha è true
	if internalDataCluster.Ha {
		yamlGroupVars["cluster_vip_fqdn"] = fmt.Sprintf("%s-vip.%s", internalDataCluster.Name, internalDataCluster.Name)
		yamlGroupVars["cluster_vip_address"] = internalDataCluster.HaIp
	}

	// Creazione del documento YAML (i trattini all'inizio)
	yamlData := []byte("---\n")

	// Conversione dell'inventario in YAML
	data, err := yaml.Marshal(yamlGroupVars)
	if err != nil {
		log.Fatalf("Failed to marshal YAML: %v", err)
	}

	// Aggiunta dei dati YAML al documento
	yamlData = append(yamlData, data...)

	// Scrittura dei dati YAML nel file di output
	writeToFile(outputFile, string(yamlData))
}
