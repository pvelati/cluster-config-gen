SHELL := /usr/bin/env bash

.PHONY: run
run:
	go run main.go

.PHONY: clear
clear:
	find ansible/. -type f -name "inventory.yaml" -prune -exec rm -rf {} \;
	find terraform/. -type f -name "*_resources.tf" -prune -exec rm -rf {} \;

