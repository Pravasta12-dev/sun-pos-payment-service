# Environment Setup Guide

## 🚀 Quick Start

### 1. Setup Environment untuk Pertama Kali
```bash
make env-setup
```
Command ini akan:
- Copy `.env.example` ke `.env`
- Tidak akan overwrite jika `.env` sudah ada

### 2. Generate Encryption Key
```bash
make env-generate-key
```
Command ini akan generate encryption key baru yang bisa dipakai untuk `ENCRYPTION_SECRET_KEY`

### 3. Refresh Environment (Hati-hati!)
```bash
make env-refresh
```
Command ini akan:
- Overwrite `.env` yang ada dengan `.env.example`
- Meminta konfirmasi dulu sebelum overwrite

## 📝 Manual Setup

Jika tidak ingin pakai Makefile, bisa manual:

```bash
# Copy .env.example ke .env
cp .env.example .env

# Generate encryption key
openssl rand -hex 32

# Edit .env file
nano .env
# atau
vim .env
```

## ⚙️ Configuration Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_ENV` | Environment mode | `development` atau `production` |
| `APP_PORT` | Application port | `8080` |
| `MIDTRANS_ENV` | Midtrans environment | `sandbox` atau `production` |
| `ENCRYPTION_SECRET_KEY` | Secret key untuk enkripsi | Generate dengan `make env-generate-key` |
| `DATABASE_HOST` | Database host | `localhost` |
| `DATABASE_PORT` | Database port | `5432` |
| `DATABASE_NAME` | Database name | `sun-pos-payment` |
| `DATABASE_USER` | Database user | `postgres` |
| `DATABASE_PASSWORD` | Database password | Your secure password |

## 🔐 Security Notes

1. **Jangan commit `.env` ke git!** (sudah ada di `.gitignore`)
2. Selalu generate `ENCRYPTION_SECRET_KEY` yang baru untuk production
3. Gunakan password yang kuat untuk database
4. Untuk production, gunakan `MIDTRANS_ENV=production`

## 🛠️ Available Make Commands

```bash
# Environment
make env-setup          # Setup .env dari .env.example (pertama kali)
make env-refresh        # Refresh .env (overwrite)
make env-generate-key   # Generate encryption key

# Database
make db-create          # Create database
make db-drop            # Drop database

# Migrations
make migrate-up         # Run migrations
make migrate-down       # Rollback migrations
make migrate-fresh      # Drop all tables and re-migrate
make migrate-create     # Create new migration

# Run
make run                # Run application
```

## 📖 Workflow Example

```bash
# 1. Clone repository
git clone <repository-url>
cd sun-pos-payment-service

# 2. Setup environment
make env-setup

# 3. Generate encryption key
make env-generate-key
# Copy output dan paste ke ENCRYPTION_SECRET_KEY di .env

# 4. Edit .env file
nano .env
# Update DATABASE_PASSWORD dan nilai lainnya

# 5. Create database
make db-create

# 6. Run migrations
make migrate-up

# 7. Run application
make run
```

## 🆘 Troubleshooting

### Error: .env file not found
```bash
make env-setup
```

### Error: database connection failed
- Check `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_USER`, `DATABASE_PASSWORD` di `.env`
- Pastikan PostgreSQL sudah running
- Pastikan database sudah dibuat dengan `make db-create`

### Error: migration failed
```bash
# Reset migrations
make migrate-fresh
```
