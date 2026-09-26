# ARCHIVED

This repository was archived on **2026-09-26**. It is preserved for
historical and reference purposes. There is no active development, no
maintenance, and no release roadmap. Issues and pull requests are closed.

## Historical purpose

`cluster-config-gen` was a small, personal "configuration compiler" for
Proxmox-based Kubernetes clusters. From a single YAML cluster definition it
generated:

- Terraform HCL (`proxmox_vm_qemu` resources plus OPNsense/Unbound DNS
  overrides) for the VM infrastructure,
- an Ansible inventory and per-cluster `group_vars` for host configuration,
- bespoke Ansible playbooks/roles installing K3s or RKE2 via the upstream
  install scripts, with a small set of local manifests (kube-vip, canal,
  coredns, ingress-nginx, system upgrade plans),
- a sops/AGE workflow for encrypting Terraform `tfvars` and S3 backend files.

The intended workflow was:

```text
config.yaml  ->  cluster-config-gen (Go)  ->  terraform/  +  ansible/
```

It was written as infrastructure automation for a single operator's lab
network (specific bridge, storage, template VM, IP layout, and VMID
conventions), not as a general-purpose tool.

## Why active development stopped

A 2026-09-26 audit (see the decision gate in the project history) concluded
that the project no longer justifies maintenance:

1. **It is glue around existing tooling.** Every generated artifact is a
   1:1 projection of inputs that Terraform/OpenTofu (with the Proxmox
   provider) and Ansible express natively. The generator adds no
   domain-specific computation beyond deriving names, IPs, and VMIDs from
   personal conventions.
2. **The unique value is personal convention.** VMID = `base_vmid + last IP
   octet`, fixed `k8s-<cluster>-master-N`/`-worker-N` naming, hard-coded
   template (`deb12-template`), storage (`local-zfs`), bridge (`vmbr20`),
   disk size (`16G`), `/24` prefixes, and hard-coded service-VIP addresses
   inside the Go source. These belong in environment-specific
   Terraform/Ansible configuration, not in a general tool.
3. **The generated output was broken at the time of archival.**
   - The shipped `config.yaml.example` makes the binary **panic** (nil
     pointer dereference in `validation.go` `overlaps()` when an IP range
     calculation fails, e.g. `last_octet: 254` with 3 nodes).
   - The generated Ansible inventory emits the key `ansiblehost` instead of
     `ansible_host` (missing YAML tag), so Ansible never used the generated
     IPs.
   - All four playbooks fail `ansible-lint` syntax checks (`cluster_name`
     undefined); the Ansible `Makefile` targets invoke playbook files as
     executables; `ansible.cfg` sets a wrong `playbook_dir` when run from
     the repository root; the `k8s_distro` variable used by
     `cluster_setup.yaml` is never defined.
   - `HaFqdn` was computed as `<domain>.<domain>` (duplicated domain) and is
     never consumed; `vip.services`, `SrvHa`, and `SrvHaIp` are dead
     configuration/fields.
   - Validation sets the HA control-plane count to 4 while generation sets
     it to 3; there is no duplicate-hostname or VIP-collision check, and IP
     allocation is last-octet string arithmetic.
4. **The K3s/RKE2 roles reimplement what maintained upstream roles already
   do.** The local roles wrap `https://get.k3s.io` / `https://get.rke2.io`
   script installs plus a handful of manifests — exactly the territory of
   maintained community roles.
5. **No evidence of external adoption.** Zero stars, zero forks, no open
   issues or pull requests, no releases, no tagged versions, no automated
   tests.

The `go build` still succeeds (Go >= 1.22), so the repository remains
buildable for reference; the known defect above means the shipped example
configuration will not run cleanly.

## Recommended maintained alternatives

For new projects, compose maintained components directly instead of this
generator:

| Concern | Maintained project |
|---|---|
| Proxmox VMs from Terraform/OpenTofu | [`bpg/terraform-provider-proxmox`](https://github.com/bpg/terraform-provider-proxmox) |
| K3s provisioning with Ansible | [`k3s-io/k3s-ansible`](https://github.com/k3s-io/k3s-ansible) |
| RKE2 provisioning with Ansible | [`rancher/rke2-ansible`](https://github.com/rancher/rke2-ansible) (also `ranchergovernment/rke2-ansible`) |
| Combined Terraform/OpenTofu + Ansible examples | Rancher infrastructure automation examples |
| Secret encryption for Terraform state/tfvars | [sops](https://github.com/getsops/sops) with AGE, as used here |

A single YAML topology file can be replaced by Terraform variables plus an
Ansible inventory; the naming/IP/VMID conventions this tool encoded are
simplest to keep as explicit values in that configuration.

## Status

- Repository: **archived, read-only, public**.
- History: preserved, not rewritten.
- Generated examples and the original implementation: preserved.
- No further feature work is planned.
