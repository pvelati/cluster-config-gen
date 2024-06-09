SHELL := /usr/bin/env bash

.PHONY: run
run:
	go build
	./cluster-config-gen

.PHONY: clear
clear:
	rm cluster-config-gen
	rm ansible/inventory.yaml
	rm terraform/*_resources.tf
