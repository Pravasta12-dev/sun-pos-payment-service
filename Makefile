include .env
export

# =====================
# ENVIRONMENT SETUP
# =====================

env-setup:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ .env file created from .env.example"; \
		echo "⚠️  Please update the values in .env file"; \
	else \
		echo "⚠️  .env file already exists. Use 'make env-refresh' to recreate it."; \
	fi

env-refresh:
	@echo "⚠️  This will overwrite your current .env file!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		cp .env.example .env; \
		echo "✅ .env file refreshed from .env.example"; \
		echo "⚠️  Please update the values in .env file"; \
	else \
		echo "❌ Operation cancelled"; \
	fi

env-generate-key:
	@echo "🔑 Generating encryption key..."
	@openssl rand -hex 32

# =====================
# DATABASE
# =====================

db-create:
	psql -h $(DATABASE_HOST) -p $(DATABASE_PORT) -U $(DATABASE_USER) -c "CREATE DATABASE $(DATABASE_NAME);"

db-drop:
	psql -h $(DATABASE_HOST) -p $(DATABASE_PORT) -U $(DATABASE_USER) -c "DROP DATABASE IF EXISTS $(DATABASE_NAME);"

# =====================
# MIGRATIONS
# =====================

migrate-up:
	migrate -database "$(DATABASE_URL)" -path database/migrations up

migrate-down:
	migrate -database "$(DATABASE_URL)" -path database/migrations down
migrate-force:
	migrate -database "$(DATABASE_URL)" -path database/migrations force $(version)

migrate-create:
	migrate create -ext sql -dir database/migrations $(name)

migrate-fresh:
	migrate -database "$(DATABASE_URL)" -path database/migrations drop
	migrate -database "$(DATABASE_URL)" -path database/migrations up
# =====================
# RUN SERVER
# =====================

run:
	go run main.go
