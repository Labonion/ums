copy_env:
	@cp .env.example .env
	@echo ".env.example copied to .env"
copy_ai_env:
	@cp ai/.env.example ai/.env
	@echo "ai/.env.example copied to ai/.env"

run_server:
	docker compose up --build

help:
	@echo "Available targets:"
	@echo "  make copy_env     - Copy env.example to .env"
	@echo "  make run_server          - Run the server"
	@echo "  make help          - Display this help message"


.DEFAULT_GOAL := help