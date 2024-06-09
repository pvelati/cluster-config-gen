provider "opnsense" {
  uri        = var.opnsense.uri
  api_key    = var.opnsense.api_key
  api_secret = var.opnsense.api_secret
}
provider "proxmox" {
  pm_api_url          = var.proxmox.api_url
  pm_api_token_id     = var.proxmox.api_token_id
  pm_api_token_secret = var.proxmox.api_token_secret
  pm_tls_insecure     = false
}
