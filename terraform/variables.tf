variable "opnsense" {
  type = object({
    uri        = string
    api_key    = string
    api_secret = string
  })
  sensitive = true
}
variable "proxmox" {
  type = object({
    api_url          = string
    api_token_id     = string
    api_token_secret = string
  })
  sensitive = true
}

variable "proxmox_node_name" {
  description = "name of the proxmox node. nodes are found in the gui under Datacenter"
}

variable "clone_type" {
  default = "full"
}

variable "qemu-guest-agent" {
  default = "enabled"
}

variable "template_name" {
  default = "deb12-template"
}

variable "master_cores" {
  default = 1
}

variable "worker_cores" {
  default = 2
}

variable "master_memory" {
  default = 1024
}

variable "worker_memory" {
  default = 2048
}

variable "bridge_name" {
  default = "vmbr20"
}

variable "storage" {
  description = "the type of storage that is backing the vm images"
  default     = "local-zfs"
}

variable "root_disk_size" {
  default = "16G"
}

variable "nameserver" {
}

variable "user" {
  default = "debian"
}

variable "sshkeys" {
}

variable "private_key" {
  description = "the ssh-key used by the provisioner to ssh into the node to finish the setup"
}
