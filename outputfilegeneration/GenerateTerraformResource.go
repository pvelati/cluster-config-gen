package outputfilegeneration

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/pvelati/cluster-config-gen/types"
	"github.com/zclconf/go-cty/cty"
)

// GenerateTerraformResource genera il file di risorse Terraform in formato HCL per un cluster specifico.
func GenerateTerraformResource(outputFile string, internalDataCluster types.InternalDataCluster) {
	// Crea il file HCL
	f := hclwrite.NewEmptyFile()

	// Ottieni il corpo principale del file
	body := f.Body()

	// Aggiungi le risorse Terraform per i master
	for _, master := range internalDataCluster.Masters {
		block := body.AppendNewBlock("resource", []string{"proxmox_vm_qemu", master.TerraformResourceName})
		blockBody := block.Body()
		blockBody.SetAttributeValue("name", cty.StringVal(master.ProxmoxVmName))
		blockBody.SetAttributeValue("desc", cty.StringVal(master.ProxmoxVmDescription))
		blockBody.SetAttributeValue("vmid", cty.NumberIntVal(int64(master.ProxmoxVMID)))
		blockBody.SetAttributeValue("full_clone", cty.BoolVal(true))
		blockBody.SetAttributeValue("agent", cty.NumberIntVal(1))
		blockBody.SetAttributeValue("clone", cty.StringVal("deb12-template"))
		blockBody.SetAttributeValue("cores", cty.NumberIntVal(1))
		blockBody.SetAttributeValue("memory", cty.NumberIntVal(1024))
		blockBody.SetAttributeValue("sockets", cty.NumberIntVal(1))
		blockBody.SetAttributeValue("cpu", cty.StringVal("host"))
		blockBody.SetAttributeValue("scsihw", cty.StringVal("virtio-scsi-single"))
		blockBody.SetAttributeValue("tags", cty.StringVal("debian;"+strings.Join(master.ProxmoxVmTags, ";")))
		// blockBody.SetAttributeValue("", cty.)
	}

	// manca la parte per i worker, ma prima voglio fare quella dei master, se va bene poi copincollo

	// Scrittura dei dati YAML nel file di output
	writeToFile(outputFile, string(f.Bytes()))
}
