SHELL := /usr/bin/env bash

.PHONY: init
init:
	terraform init -backend-config=config.s3.tfbackend 

.PHONY: clear
clear:
	find . -type f -name ".terraform.lock.hcl" -prune -exec rm -rf {} \;
	rm -rf .terraform

.PHONY: validate
validate:
	terraform fmt && terraform validate

.PHONY: plan
plan:
	terraform plan

.PHONY: apply
apply:
	terraform apply -auto-approve

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

