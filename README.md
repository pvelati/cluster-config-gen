# cluster-config-gen

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

