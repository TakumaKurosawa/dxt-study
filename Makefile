.DEFAULT_GOAL := help

##### SETUP #####

.PHONY: setup
setup: ## Install all dependencies ## make setup
	pnpm install
	go work sync

	@if ! command -v dxt >/dev/null 2>&1; then \
		echo "🔧 dxt command not found. Installing..."; \
		pnpm i -g @anthropic-ai/dxt; \
		echo "✅ dxt installed successfully"; \
	else \
		echo "✅ dxt command is already installed"; \
	fi

	@if ! command -v gitleaks >/dev/null 2>&1; then \
		echo "🔒 gitleaks command not found. Installing..."; \
		go install github.com/zricethezav/gitleaks/v8@latest; \
		echo "✅ gitleaks installed successfully"; \
	else \
		echo "✅ gitleaks command is already installed"; \
	fi
	
	@if ! command -v lefthook >/dev/null 2>&1; then \
		echo "🪝 lefthook command not found. Installing..."; \
		pnpm i -g lefthook; \
		echo "✅ lefthook installed successfully"; \
	else \
		echo "✅ lefthook command is already installed"; \
	fi

	@echo "🔧 Setting up lefthook hooks..."
	lefthook install

##### RUN #####

.PHONY: dev/hello-world
dev/hello-world: ## Start hello-world dev server ## make dev/hello-world
	pnpm run dev:hello-world

.PHONY: dev/hello-world-binary
dev/hello-world-binary: ## Start hello-world-binary dev server ## make dev/hello-world-binary
	cd apps/hello-world-binary/server && go run main.go

##### BUILD DXT #####

.PHONY: build
build: ## Build DXT application ## make build [app=hello-world]
build: app ?= hello-world
build:
	@if [ ! -d "apps/$(app)" ]; then \
		echo "❌ Error: Application directory 'apps/$(app)' does not exist"; \
		echo ""; \
		echo "Available apps:"; \
		ls apps/ | sed 's/^/  - /'; \
		exit 1; \
	fi
	@echo "🔨 Building DXT for $(app)..."
	cd apps/$(app) && dxt pack
	@echo "✅ Build completed for $(app)"

##### CODE QUALITY #####

.PHONY: lint
lint: ## Run linter on all files ## make lint
	pnpm run lint

.PHONY: format
format: ## Format all files ## make format
	pnpm run format

##### HELP #####

.PHONY: help
help: ## Display this help screen ## make or make help
	@echo ""
	@echo "Usage: make SUB_COMMAND argument_name=argument_value"
	@echo ""
	@echo "Command list:"
	@echo ""
	@printf "\033[36m%-30s\033[0m %-50s %s\n" "[Sub command]" "[Description]" "[Example]"
	@grep -E '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | perl -pe 's%^([/a-zA-Z_-]+):.*?(##)%$$1 $$2%' | awk -F " *?## *?" '{printf "\033[36m%-30s\033[0m %-50s %s\n", $$1, $$2, $$3}'
