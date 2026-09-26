# cluster-config-gen

> ## ⚠️ ARCHIVED (2026-09-26)
>
> This project is **archived and read-only**. It is preserved for
> historical/reference purposes only — no maintenance, no releases, no new
> features. See **[ARCHIVED.md](ARCHIVED.md)** for the full rationale and for
> the maintained alternatives to use instead
> (`bpg/terraform-provider-proxmox`, `k3s-io/k3s-ansible`, `rancher/rke2-ansible`).
>
> Note: the shipped `config.yaml.example` panics the binary (nil-pointer in
> validation); the generated Ansible inventory used an invalid `ansiblehost`
> key. Do not use this tool for new infrastructure.

---

**Generate full RKE2/K3s clusters on Proxmox.**

Configuration-driven tool used to generate k3s or rke2 clusters on Proxmox.

With a single YAML configuration file, it produces ready-to-use Terraform HCL files and Ansible inventories/vars.

---

## 🚀 Features

- Supports both **RKE2** and **K3s** deployments
- Cluster definitions via a single YAML config file
- Generates:
  - Proxmox infrastructure code using **Terraform**
  - Host configuration using **Ansible**
- Works with **AGE** encryption for secrets
- Makefile-driven workflow for simplicity and automation
- Multi-cluster ready (e.g. `prod`, `dev`, etc.)

---

## 📁 Example Configuration

See config.yaml.example

---

## ⚙️ Usage

Use the makefile, check inside for the commands

---

## 🔐 Secrets and Encryption
This project uses AGE for encrypting secrets and Ansible variables.

