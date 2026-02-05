include .env
export

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
