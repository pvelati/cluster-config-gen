SHELL := /usr/bin/env bash

CLUSTER=$1

.PHONY: k3s CLUSTER=<cluster_name>
k3s:
	playbooks/k3s_cluster_setup.yaml -e cluster_name=$(CLUSTER)

.PHONY: rke2 CLUSTER=<cluster_name>
rke2:
	playbooks/rke2_cluster_setup.yaml -e cluster_name=$(CLUSTER)

.PHONY: rke2 CLUSTER=<cluster_name>
setup:
	playbooks/cluster_setup.yaml -e cluster_name=$(CLUSTER)

.PHONY: destroy CLUSTER=<cluster_name>
destroy:
	playbooks/uninstall_cluster.yaml -e cluster_name=$(CLUSTER)

.PHONY: encrypt
encrypt:
	sops --encrypt --age `cat age-key.txt | grep -oP "public key: \K(.*)"` ./terraform.tfvars.json > terraform.tfvars.json.enc
	sops --encrypt --age `cat age-key.txt | grep -oP "public key: \K(.*)"` ./config.s3.tfbackend > config.s3.tfbackend.enc

.PHONY: decrypt
decrypt:
	SOPS_AGE_KEY_FILE=age-key.txt sops --output-type json --decrypt ./terraform.tfvars.json.enc > terraform.tfvars.json 
	SOPS_AGE_KEY_FILE=age-key.txt sops --decrypt ./config.s3.tfbackend.enc > config.s3.tfbackend 

.PHONY: test
test:
	 `cat age-key.txt | grep -oP "public key: \K(.*)"`

