package outputfilegeneration

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
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

	// Aggiungi la entry del vip se definito nella config
	if internalDataCluster.Ha {
		dnsvip := body.AppendNewBlock("resource", []string{"opnsense_unbound_host_override", internalDataCluster.Name})
		dnsvipBody := dnsvip.Body()
		dnsvipBody.SetAttributeValue("enabled", cty.BoolVal(true))
		dnsvipBody.SetAttributeValue("description", cty.StringVal(internalDataCluster.Name+" VIP address"))
		dnsvipBody.SetAttributeValue("hostname", cty.StringVal("k8s-"+internalDataCluster.Name+"-vip"))
		dnsvipBody.SetAttributeValue("domain", cty.StringVal(internalDataCluster.Masters[0].Domain))
		dnsvipBody.SetAttributeValue("server", cty.StringVal(internalDataCluster.HaIp))
	}

	// Aggiungi le risorse Terraform per tutti i nodi (master e worker)
	for _, node := range append(internalDataCluster.Masters, internalDataCluster.Workers...) {
		dns := body.AppendNewBlock("resource", []string{"opnsense_unbound_host_override", node.TerraformResourceName + "_dns"})
		dnsBody := dns.Body()
		dnsBody.SetAttributeValue("enabled", cty.BoolVal(true))
		dnsBody.SetAttributeValue("description", cty.StringVal(node.ProxmoxVmDescription))
		dnsBody.SetAttributeValue("hostname", cty.StringVal(node.ProxmoxVmName))
		dnsBody.SetAttributeValue("domain", cty.StringVal(node.Domain))
		dnsBody.SetAttributeValue("server", cty.StringVal(node.IP))

		vm := body.AppendNewBlock("resource", []string{"proxmox_vm_qemu", node.TerraformResourceName})
		vmBody := vm.Body()
		vmBody.SetAttributeValue("name", cty.StringVal(node.ProxmoxVmName))
		vmBody.SetAttributeValue("desc", cty.StringVal(node.ProxmoxVmDescription))
		vmBody.SetAttributeValue("vmid", cty.NumberIntVal(int64(node.ProxmoxVMID)))
		vmBody.SetAttributeValue("full_clone", cty.BoolVal(true))
		vmBody.SetAttributeValue("agent", cty.NumberIntVal(1))
		vmBody.SetAttributeTraversal("target_node", hcl.Traversal{
			hcl.TraverseRoot{
				Name: "var",
			},
			hcl.TraverseAttr{
				Name: "proxmox_node_name",
			},
		})
		vmBody.SetAttributeValue("clone", cty.StringVal("deb12-template"))
		vmBody.SetAttributeValue("cores", cty.NumberIntVal(int64(node.Core)))
		vmBody.SetAttributeValue("memory", cty.NumberIntVal(int64(node.Memory)))
		vmBody.SetAttributeValue("sockets", cty.NumberIntVal(1))
		vmBody.SetAttributeValue("cpu", cty.StringVal("host"))
		vmBody.SetAttributeValue("scsihw", cty.StringVal("virtio-scsi-single"))
		vmBody.SetAttributeValue("tags", cty.StringVal("debian;"+strings.Join(node.ProxmoxVmTags, ";")))
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
		vmBody.SetAttributeValue("skip_ipv6", cty.BoolVal(true))
		vmBody.SetAttributeValue("agent_timeout", cty.NumberIntVal(120))
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
		vmBody.SetAttributeValue("ipconfig0", cty.StringVal("ip="+node.IP+"/24,gw="+node.Gateway))
		vmBody.SetAttributeValue("searchdomain", cty.StringVal(node.Domain))
		vmBody.SetAttributeTraversal("nameserver", hcl.Traversal{
			hcl.TraverseRoot{
				Name: "var",
			},
			hcl.TraverseAttr{
				Name: "nameserver",
			},
		})
		vmBody.SetAttributeTraversal("ciuser", hcl.Traversal{
			hcl.TraverseRoot{
				Name: "var",
			},
			hcl.TraverseAttr{
				Name: "user",
			},
		})
		vmBody.SetAttributeTraversal("sshkeys", hcl.Traversal{
			hcl.TraverseRoot{
				Name: "var",
			},
			hcl.TraverseAttr{
				Name: "sshkeys",
			},
		})
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
		vmConnectionBody.SetAttributeTraversal("host", hcl.Traversal{
			hcl.TraverseRoot{
				Name: "self",
			},
			hcl.TraverseAttr{
				Name: "default_ipv4_address",
			},
		})
		vmConnectionBody.SetAttributeTraversal("user", hcl.Traversal{
			hcl.TraverseRoot{
				Name: "var",
			},
			hcl.TraverseAttr{
				Name: "user",
			},
		})
		vmConnectionBody.SetAttributeTraversal("private_key", hcl.Traversal{
			hcl.TraverseRoot{
				Name: "var",
			},
			hcl.TraverseAttr{
				Name: "private_key",
			},
		})
		vmProvisioner := vmBody.AppendNewBlock("provisioner", []string{"remote-exec"})
		vmProvisionerBody := vmProvisioner.Body()
		vmProvisionerBody.SetAttributeValue("inline", cty.ListVal([]cty.Value{
			cty.StringVal("sudo systemd-machine-id-setup"),
			cty.StringVal("sudo shutdown -r +0"),
		}))
	}

	// Scrittura dei dati YAML nel file di output
	writeToFile(outputFile, string(f.Bytes()))
}
