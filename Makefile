SHELL := /usr/bin/env bash

.PHONY: run
run:
	go run main.go

.PHONY: clear
clear:
	find outputs/. -type f -name "inventory.yaml" -prune -exec rm -rf {} \;
	find outputs/. -type f -name "*_resources.tf" -prune -exec rm -rf {} \;

