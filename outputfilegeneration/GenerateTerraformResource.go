package outputfilegeneration

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
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
		vm := body.AppendNewBlock("resource", []string{"proxmox_vm_qemu", master.TerraformResourceName})
		vmBody := vm.Body()
		vmBody.SetAttributeValue("name", cty.StringVal(master.ProxmoxVmName))
		vmBody.SetAttributeValue("desc", cty.StringVal(master.ProxmoxVmDescription))
		vmBody.SetAttributeValue("vmid", cty.NumberIntVal(int64(master.ProxmoxVMID)))
		vmBody.SetAttributeValue("full_clone", cty.BoolVal(true))
		vmBody.SetAttributeValue("agent", cty.NumberIntVal(1))
		vmBody.SetAttributeValue("target_node", cty.StringVal("var.proxmox_node_name"))
		vmBody.SetAttributeValue("clone", cty.StringVal("deb12-template"))
		vmBody.SetAttributeValue("cores", cty.NumberIntVal(1))
		vmBody.SetAttributeValue("memory", cty.NumberIntVal(1024))
		vmBody.SetAttributeValue("sockets", cty.NumberIntVal(1))
		vmBody.SetAttributeValue("cpu", cty.StringVal("host"))
		vmBody.SetAttributeValue("scsihw", cty.StringVal("virtio-scsi-single"))
		vmBody.SetAttributeValue("tags", cty.StringVal("debian;"+strings.Join(master.ProxmoxVmTags, ";")))
		vmBody.AppendNewline()

		vmNetwork := vmBody.AppendNewBlock("network", nil)
		vmNetworkBody := vmNetwork.Body()
		vmNetworkBody.SetAttributeValue("bridge", cty.StringVal("vmbr20"))
		vmNetworkBody.SetAttributeValue("firewall", cty.BoolVal(false))
		vmNetworkBody.SetAttributeValue("model", cty.StringVal("virtio"))
		vmBody.AppendNewline()

		commentDisks := hclwrite.Tokens{
			&hclwrite.Token{
				Type:         hclsyntax.TokenComment,
				Bytes:        []byte("// This must match what you've set in the template or you will end up with 2 disks\n"),
				SpacesBefore: 0,
			},
		}
		vmBody.AppendUnstructuredTokens(commentDisks)
		vmDisks := vmBody.AppendNewBlock("disks", nil)
		vmDisksBody := vmDisks.Body()
		vmDisksScsi := vmDisksBody.AppendNewBlock("scsi", nil)
		vmDisksScsiBody := vmDisksScsi.Body()
		vmDisksScsiScsi0 := vmDisksScsiBody.AppendNewBlock("scsi0", nil)
		vmDisksScsiScsi0Body := vmDisksScsiScsi0.Body()
		vmDisksScsiScsi0Disk := vmDisksScsiScsi0Body.AppendNewBlock("disk", nil)
		vmDisksScsiScsi0DiskBody := vmDisksScsiScsi0Disk.Body()
		vmDisksScsiScsi0DiskBody.SetAttributeValue("asyncio", cty.StringVal("native"))
		vmDisksScsiScsi0DiskBody.SetAttributeValue("cache", cty.StringVal("none"))
		vmDisksScsiScsi0DiskBody.SetAttributeValue("discard", cty.BoolVal(true))
		vmDisksScsiScsi0DiskBody.SetAttributeValue("emulatessd", cty.BoolVal(true))
		vmDisksScsiScsi0DiskBody.SetAttributeValue("iothread", cty.BoolVal(true))
		vmDisksScsiScsi0DiskBody.SetAttributeValue("storage", cty.StringVal("local-zfs"))
		vmDisksScsiScsi0DiskBody.SetAttributeValue("size", cty.StringVal("16G"))
		vmDisksIde := vmDisksBody.AppendNewBlock("ide", nil)
		vmDisksIdeBody := vmDisksIde.Body()
		vmDisksIdeIde2 := vmDisksIdeBody.AppendNewBlock("ide2", nil)
		vmDisksIdeIde2Body := vmDisksIdeIde2.Body()
		vmDisksIdeIde2Cloudinit := vmDisksIdeIde2Body.AppendNewBlock("cloudinit", nil)
		vmDisksIdeIde2CloudinitBody := vmDisksIdeIde2Cloudinit.Body()
		vmDisksIdeIde2CloudinitBody.SetAttributeValue("storage", cty.StringVal("local-zfs"))
		vmBody.AppendNewline()

		commentCloudinit := hclwrite.Tokens{
			&hclwrite.Token{
				Type:         hclsyntax.TokenComment,
				Bytes:        []byte("// CLOUD-INIT\n"),
				SpacesBefore: 0,
			},
		}
		vmBody.AppendUnstructuredTokens(commentCloudinit)
		vmBody.SetAttributeValue("os_type", cty.StringVal("cloud-init"))
		vmBody.SetAttributeValue("ipconfig0", cty.StringVal("ip="+master.IP+"/24,gw=${var.gateway}"))
		vmBody.SetAttributeValue("nameserver", cty.StringVal("var.nameserver"))
		vmBody.SetAttributeValue("ciuser", cty.StringVal("var.user"))
		vmBody.SetAttributeValue("sshkeys", cty.StringVal("var.sshkeys"))
		vmBody.AppendNewline()

		commentProvisioner := hclwrite.Tokens{
			&hclwrite.Token{
				Type:         hclsyntax.TokenComment,
				Bytes:        []byte("// PROVISIONER\n"),
				SpacesBefore: 0,
			},
		}
		vmBody.AppendUnstructuredTokens(commentProvisioner)
		vmConnection := vmBody.AppendNewBlock("connection", nil)
		vmConnectionBody := vmConnection.Body()
		vmConnectionBody.SetAttributeValue("type", cty.StringVal("ssh"))
		vmConnectionBody.SetAttributeValue("host", cty.StringVal("self.default_ipv4_address"))
		vmConnectionBody.SetAttributeValue("user", cty.StringVal("var.user"))
		vmConnectionBody.SetAttributeValue("private_key", cty.StringVal("var.private_key"))
		vmProvisioner := vmBody.AppendNewBlock("provisioner", []string{"remote-exec"})
		vmProvisionerBody := vmProvisioner.Body()
		vmProvisionerBody.SetAttributeValue("inline", cty.ListVal([]cty.Value{
			cty.StringVal("sudo systemd-machine-id-setup"),
			cty.StringVal("sudo shutdown -r +0"),
		}))
	}

	// manca la parte per i worker, ma prima voglio fare quella dei master, se va bene poi copincollo
	// oppure ciclare usando il blocco precendente con if master/worker e mettergli dentro i valori?

	// Scrittura dei dati YAML nel file di output
	writeToFile(outputFile, string(f.Bytes()))
}
