#!/bin/bash
# Database backup script for Opslane
# Usage: ./scripts/backup.sh [backup_name]

set -e

# Configuration
DB_NAME="${OPSLANE_DB_NAME:-opslane}"
DB_USER="${OPSLANE_DB_USER:-opslane}"
DB_HOST="${OPSLANE_DB_HOST:-localhost}"
DB_PORT="${OPSLANE_DB_PORT:-5432}"
BACKUP_DIR="${OPSLANE_BACKUP_DIR:-./backups}"

# Get timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="${1:-opslane_${TIMESTAMP}}"
BACKUP_FILE="${BACKUP_DIR}/${BACKUP_NAME}.sql"

# Create backup directory if it doesn't exist
mkdir -p "${BACKUP_DIR}"

echo "Creating database backup: ${BACKUP_NAME}"
echo "Database: ${DB_NAME}"
echo "Host: ${DB_HOST}:${DB_PORT}"

# Check if we're in Docker
if [ -f /.dockerenv ] || grep -q docker /proc/1/cgroup 2>/dev/null; then
    # We're in Docker, use the db hostname
    pg_dump -h db -U "${DB_USER}" "${DB_NAME}" -F p -f "${BACKUP_FILE}"
else
    # We're on host, use localhost
    pg_dump -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" "${DB_NAME}" -F p -f "${BACKUP_FILE}"
fi

# Get file size
SIZE=$(du -h "${BACKUP_FILE}" | cut -f1)

echo ""
echo "✓ Backup created: ${BACKUP_FILE} (${SIZE})"

# Keep only last 7 backups
cd "${BACKUP_DIR}"
ls -t opslane_*.sql 2>/dev/null | tail -n +8 | xargs -r rm -f

echo "✓ Old backups cleaned up (keeping last 7)"