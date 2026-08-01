.PHONY: help lint cluster-up tofu-fmt tofu-plan

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: ## Run pre-commit hooks for all files
	pre-commit run --all-files

cluster-up: ## Spin up a local K3s cluster using Ansible
	cd ansible && ansible-playbook -i inventory.ini k3s-installation/playbook.yaml

tofu-fmt: ## Format OpenTofu code
	cd infrastructure && tofu fmt -recursive

tofu-plan: ## Run OpenTofu plan
	cd infrastructure && tofu init && tofu plan
