terraform {
  required_providers {
    opnsense = {
      source  = "browningluke/opnsense"
      version = "0.10.1"
    }
    proxmox = {
      source  = "registry.terraform.io/telmate/proxmox"
      version = "3.0.1-rc2"
    }
  }
}
